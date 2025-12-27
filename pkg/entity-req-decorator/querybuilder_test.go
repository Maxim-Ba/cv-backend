package entityreqdecorator

import (
	"strings"
	"testing"
)

func testFieldValidator(field string) bool {
	allowedFields := map[string]bool{
		"name":    true,
		"age":     true,
		"email":   true,
		"id":      true,
		"status":  true,
		"created": true,
	}
	return allowedFields[field]
}
func TestAddFilter(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		predicate     SQLGenerator
		fieldValidator func(string) bool
		wantCondition string
		wantParams    []interface{}
	}{
		{
			name:  "PredicateEQ with valid field",
			field: "name",
			predicate: &PredicateEQ{
				Predicate: Predicate{
					Value: "John",
				},
			},
			fieldValidator: testFieldValidator,
			wantCondition: "name = $1",
			wantParams:    []interface{}{"John"},
		},
		{
			name:  "PredicateLike with valid field",
			field: "email",
			predicate: &PredicateLike{
				Predicate: Predicate{
					Value: "test",
				},
			},
			fieldValidator: testFieldValidator,
			wantCondition: "email LIKE $1",
			wantParams:    []interface{}{"%test%"},
		},
		{
			name:  "PredicateANF with multiple conditions",
			field: "age",
			predicate: &PredicateANF{
				Predicate: Predicate{
					InnerPredicate: []SQLGenerator{
						&PredicateGT{
							Predicate: Predicate{Value: "18"},
						},
						&PredicateLT{
							Predicate: Predicate{Value: "65"},
						},
					},
				},
			},
			fieldValidator: testFieldValidator,
			wantCondition: "(age > $1 AND age < $2)",
			wantParams:    []interface{}{"18", "65"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			qb.AddFilter(tt.field, tt.predicate, tt.fieldValidator)

			// Проверяем условие WHERE
			whereClause := qb.BuildWhereClause()
			if tt.wantCondition != "" {
				expectedWhere := " WHERE " + tt.wantCondition
				if whereClause != expectedWhere {
					t.Errorf("BuildWhereClause() = %v, want %v", whereClause, expectedWhere)
				}
			} else if whereClause != "" {
				t.Errorf("BuildWhereClause() = %v, want empty", whereClause)
			}

			// Проверяем параметры
			params := qb.GetParams()
			if len(params) != len(tt.wantParams) {
				t.Errorf("GetParams() length = %v, want %v", len(params), len(tt.wantParams))
			}
			for i, param := range params {
				if param != tt.wantParams[i] {
					t.Errorf("GetParams()[%d] = %v, want %v", i, param, tt.wantParams[i])
				}
			}
		})
	}
}

func TestAddSort(t *testing.T) {
	tests := []struct {
		name          string
		sorts         []SortBy
		fieldValidator func(string) bool
		wantOrderBy   string
	}{
		{
			name: "single valid sort",
			sorts: []SortBy{
				{Field: "name", Order: "ASC"},
			},
			fieldValidator: testFieldValidator,
			wantOrderBy:   " ORDER BY name ASC",
		},
		{
			name: "multiple valid sorts",
			sorts: []SortBy{
				{Field: "name", Order: "ASC"},
				{Field: "id", Order: "DESC"},
			},
			fieldValidator: testFieldValidator,
			wantOrderBy:   " ORDER BY name ASC, id DESC",
		},
		{
			name: "mixed valid and invalid sorts",
			sorts: []SortBy{
				{Field: "name", Order: "ASC"},
				{Field: "invalid_field", Order: "DESC"},
				{Field: "email", Order: "ASC"},
			},
			fieldValidator: testFieldValidator,
			wantOrderBy:   " ORDER BY name ASC, email ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			qb.AddSort(tt.sorts, tt.fieldValidator)

			if qb.OrderByClause != tt.wantOrderBy {
				t.Errorf("OrderByClause = %v, want %v", qb.OrderByClause, tt.wantOrderBy)
			}
		})
	}
}

func TestAddPagination(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		size        int
		wantLimit   string
		wantOffset  string
		wantParams  []interface{}
	}{
		{
			name:       "first page with limit",
			page:       1,
			size:       10,
			wantLimit:  " LIMIT $1",
			wantOffset: "",
			wantParams: []interface{}{10},
		},
		{
			name:       "second page with limit and offset",
			page:       2,
			size:       20,
			wantLimit:  " LIMIT $1",
			wantOffset: " OFFSET $2",
			wantParams: []interface{}{20, 20}, // (2-1)*20 = 20
		},
		{
			name:       "page with zero size",
			page:       3,
			size:       0,
			wantLimit:  "",
			wantOffset: "",
			wantParams: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			qb.AddPagination(tt.page, tt.size)

			if qb.LimitClause != tt.wantLimit {
				t.Errorf("LimitClause = %v, want %v", qb.LimitClause, tt.wantLimit)
			}
			if qb.OffsetClause != tt.wantOffset {
				t.Errorf("OffsetClause = %v, want %v", qb.OffsetClause, tt.wantOffset)
			}
			
			params := qb.GetParams()
			if len(params) != len(tt.wantParams) {
				t.Errorf("GetParams() length = %v, want %v", len(params), len(tt.wantParams))
			}
			for i, param := range params {
				if param != tt.wantParams[i] {
					t.Errorf("GetParams()[%d] = %v, want %v", i, param, tt.wantParams[i])
				}
			}
		})
	}
}


func TestBuildWhereClause(t *testing.T) {
	tests := []struct {
		name           string
		setupBuilder   func(*QueryBuilder)
		wantWhereClause string
	}{
		{
			name: "no conditions",
			setupBuilder: func(qb *QueryBuilder) {
				// Не добавляем условий
			},
			wantWhereClause: "",
		},
		{
			name: "single condition",
			setupBuilder: func(qb *QueryBuilder) {
				qb.AddFilter("name", &PredicateEQ{
					Predicate: Predicate{Value: "John"},
				}, testFieldValidator)
			},
			wantWhereClause: " WHERE name = $1",
		},
		{
			name: "multiple conditions",
			setupBuilder: func(qb *QueryBuilder) {
				qb.AddFilter("name", &PredicateEQ{
					Predicate: Predicate{Value: "John"},
				}, testFieldValidator)
				qb.AddFilter("age", &PredicateGT{
					Predicate: Predicate{Value: "18"},
				}, testFieldValidator)
			},
			wantWhereClause: " WHERE name = $1 AND age > $2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			tt.setupBuilder(qb)
			
			whereClause := qb.BuildWhereClause()
			if whereClause != tt.wantWhereClause {
				t.Errorf("BuildWhereClause() = %v, want %v", whereClause, tt.wantWhereClause)
			}
		})
	}
}


func TestBuildSelectQuery(t *testing.T) {
	tests := []struct {
		name         string
		baseQuery    string
		setupBuilder func(*QueryBuilder)
		wantQuery    string
	}{
		{
			name:      "base query only",
			baseQuery: "SELECT * FROM users",
			setupBuilder: func(qb *QueryBuilder) {
				// Не добавляем условий
			},
			wantQuery: "SELECT * FROM users",
		},
		{
			name:      "query with where and order by",
			baseQuery: "SELECT id, name FROM users",
			setupBuilder: func(qb *QueryBuilder) {
				qb.AddFilter("name", &PredicateEQ{
					Predicate: Predicate{Value: "John"},
				}, testFieldValidator)
				qb.AddSort([]SortBy{
					{Field: "id", Order: "DESC"},
				}, testFieldValidator)
			},
			wantQuery: "SELECT id, name FROM users WHERE name = $1 ORDER BY id DESC",
		},
		{
			name:      "full query with all clauses",
			baseQuery: "SELECT * FROM products",
			setupBuilder: func(qb *QueryBuilder) {
				qb.AddFilter("status", &PredicateEQ{
					Predicate: Predicate{Value: "active"},
				}, testFieldValidator)
				qb.AddSort([]SortBy{
					{Field: "created", Order: "DESC"},
					{Field: "name", Order: "ASC"},
				}, testFieldValidator)
				qb.AddPagination(2, 10)
			},
			wantQuery: "SELECT * FROM products WHERE status = $1 ORDER BY created DESC, name ASC LIMIT $2 OFFSET $3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			tt.setupBuilder(qb)
			
			query := qb.BuildSelectQuery(tt.baseQuery)
			if query != tt.wantQuery {
				t.Errorf("BuildSelectQuery() = %v, want %v", query, tt.wantQuery)
			}
		})
	}
}

func TestBuildCountQuery(t *testing.T) {
	tests := []struct {
		name         string
		baseQuery    string
		setupBuilder func(*QueryBuilder)
		wantQuery    string
	}{
		{
			name:      "base query without conditions",
			baseQuery: "SELECT * FROM users",
			setupBuilder: func(qb *QueryBuilder) {
				// Не добавляем условий
			},
			wantQuery: "SELECT COUNT(*) FROM (SELECT * FROM users) as subquery",
		},
		{
			name:      "query with where conditions",
			baseQuery: "SELECT id, name FROM users",
			setupBuilder: func(qb *QueryBuilder) {
				qb.AddFilter("name", &PredicateLike{
					Predicate: Predicate{Value: "John"},
				}, testFieldValidator)
				qb.AddFilter("age", &PredicateGT{
					Predicate: Predicate{Value: "18"},
				}, testFieldValidator)
			},
			wantQuery: "SELECT COUNT(*) FROM (SELECT id, name FROM users WHERE name LIKE $1 AND age > $2) as subquery",
		},
		{
			name:      "query with where but no pagination params",
			baseQuery: "SELECT * FROM products",
			setupBuilder: func(qb *QueryBuilder) {
				qb.AddFilter("status", &PredicateEQ{
					Predicate: Predicate{Value: "active"},
				}, testFieldValidator)
				qb.AddPagination(2, 10) // Добавляем пагинацию, но она не должна влиять на COUNT
			},
			wantQuery: "SELECT COUNT(*) FROM (SELECT * FROM products WHERE status = $1) as subquery",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			tt.setupBuilder(qb)
			
			query := qb.BuildCountQuery(tt.baseQuery)
			if query != tt.wantQuery {
				t.Errorf("BuildCountQuery() = %v, want %v", query, tt.wantQuery)
			}
			
			// Проверяем, что COUNT запрос не содержит LIMIT и OFFSET
			if strings.Contains(strings.ToUpper(query), "LIMIT") {
				t.Errorf("BuildCountQuery() contains LIMIT, but shouldn't")
			}
			if strings.Contains(strings.ToUpper(query), "OFFSET") {
				t.Errorf("BuildCountQuery() contains OFFSET, but shouldn't")
			}
		})
	}
}

func TestGetCountParams(t *testing.T) {
	tests := []struct {
		name         string
		setupBuilder func(*QueryBuilder)
		wantParams   []interface{}
	}{
		{
			name: "only where conditions",
			setupBuilder: func(qb *QueryBuilder) {
				qb.AddFilter("name", &PredicateEQ{
					Predicate: Predicate{Value: "John"},
				}, testFieldValidator)
				qb.AddFilter("age", &PredicateGT{
					Predicate: Predicate{Value: "18"},
				}, testFieldValidator)
			},
			wantParams: []interface{}{"John", "18"},
		},
		{
			name: "where conditions with pagination",
			setupBuilder: func(qb *QueryBuilder) {
				qb.AddFilter("status", &PredicateEQ{
					Predicate: Predicate{Value: "active"},
				}, testFieldValidator)
				qb.AddPagination(2, 10)
			},
			wantParams: []interface{}{"active"},
		},
		{
			name: "only pagination without where",
			setupBuilder: func(qb *QueryBuilder) {
				qb.AddPagination(3, 20)
			},
			wantParams: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewQueryBuilder()
			tt.setupBuilder(qb)
			
			params := qb.GetCountParams()
			
			if len(params) != len(tt.wantParams) {
				t.Errorf("GetCountParams() length = %v, want %v", len(params), len(tt.wantParams))
			}
			
			for i, param := range params {
				if param != tt.wantParams[i] {
					t.Errorf("GetCountParams()[%d] = %v, want %v", i, param, tt.wantParams[i])
				}
			}
		})
	}
}
