// querybuilder.go
package entityreqdecorator

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	WhereConditions []string
	OrderByClause   string // " ORDER BY name ASC, id DESC"
	LimitClause     string //" LIMIT $5"
	OffsetClause    string //" OFFSET $6"
	Params          []interface{}
	paramCounter    int // Счетчик параметров (начинается с 1 для $1, $2, ...)
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		WhereConditions: make([]string, 0),
		Params:          make([]interface{}, 0),
		paramCounter:    1,
	}
}

// AddFilter добавляет условие фильтрации в WHERE часть запроса.
// field: имя поля таблицы
// predicate: предикат (например, PredicateEQ, PredicateGT и т.д.)
// fieldValidator: функция валидации имени поля (возвращает true, если поле допустимо)
// Пример для predicate типа PredicateEQ{Value: "test"}:
//   Добавляет условие: "name = $1"
//   Добавляет параметр: "test" в qb.Params
//   Увеличивает paramCounter: с 1 на 2
func (qb *QueryBuilder) AddFilter(field string, predicate SQLGenerator, fieldValidator func(string) bool) {
	if !fieldValidator(field) {
		return
	}

	sqlPart := predicate.GenerateSQL()
	if sqlPart == "" {
		return
	}

	switch p := predicate.(type) {
	case *PredicateEQ:
		qb.WhereConditions = append(qb.WhereConditions, fmt.Sprintf("%s = $%d", field, qb.paramCounter))
		qb.Params = append(qb.Params, p.Value)
		qb.paramCounter++

	case *PredicateLike:
		qb.WhereConditions = append(qb.WhereConditions, fmt.Sprintf("%s LIKE $%d", field, qb.paramCounter))
		qb.Params = append(qb.Params, "%"+p.Value+"%")
		qb.paramCounter++

	case *PredicateGT:
		qb.WhereConditions = append(qb.WhereConditions, fmt.Sprintf("%s > $%d", field, qb.paramCounter))
		qb.Params = append(qb.Params, p.Value)
		qb.paramCounter++

	case *PredicateLT:
		qb.WhereConditions = append(qb.WhereConditions, fmt.Sprintf("%s < $%d", field, qb.paramCounter))
		qb.Params = append(qb.Params, p.Value)
		qb.paramCounter++

	case *PredicateGTE:
		qb.WhereConditions = append(qb.WhereConditions, fmt.Sprintf("%s >= $%d", field, qb.paramCounter))
		qb.Params = append(qb.Params, p.Value)
		qb.paramCounter++

	case *PredicateLTE:
		qb.WhereConditions = append(qb.WhereConditions, fmt.Sprintf("%s <= $%d", field, qb.paramCounter))
		qb.Params = append(qb.Params, p.Value)
		qb.paramCounter++

	case *PredicateNE:
		qb.WhereConditions = append(qb.WhereConditions, fmt.Sprintf("%s != $%d", field, qb.paramCounter))
		qb.Params = append(qb.Params, p.Value)
		qb.paramCounter++

	case *PredicateANF:
		var anfConditions []string
		for _, innerPred := range p.InnerPredicate {
			switch ip := innerPred.(type) {
			case *PredicateGT:
				anfConditions = append(anfConditions, fmt.Sprintf("%s > $%d", field, qb.paramCounter))
				qb.Params = append(qb.Params, ip.Value)
				qb.paramCounter++

			case *PredicateLT:
				anfConditions = append(anfConditions, fmt.Sprintf("%s < $%d", field, qb.paramCounter))
				qb.Params = append(qb.Params, ip.Value)
				qb.paramCounter++

			case *PredicateGTE:
				anfConditions = append(anfConditions, fmt.Sprintf("%s >= $%d", field, qb.paramCounter))
				qb.Params = append(qb.Params, ip.Value)
				qb.paramCounter++

			case *PredicateLTE:
				anfConditions = append(anfConditions, fmt.Sprintf("%s <= $%d", field, qb.paramCounter))
				qb.Params = append(qb.Params, ip.Value)
				qb.paramCounter++

			case *PredicateEQ:
				anfConditions = append(anfConditions, fmt.Sprintf("%s = $%d", field, qb.paramCounter))
				qb.Params = append(qb.Params, ip.Value)
				qb.paramCounter++

			case *PredicateLike:
				anfConditions = append(anfConditions, fmt.Sprintf("%s LIKE $%d", field, qb.paramCounter))
				qb.Params = append(qb.Params, "%"+ip.Value+"%")
				qb.paramCounter++
			}
		}

		if len(anfConditions) > 0 {
			qb.WhereConditions = append(qb.WhereConditions, "("+strings.Join(anfConditions, " AND ")+")")
		}
	}
}

// AddSort добавляет сортировку в запрос на основе массива SortBy.
// sorts: массив структур SortBy с полем Field и порядком Order (ASC/DESC)
// fieldValidator: функция валидации имени поля
// Пример для sorts = []SortBy{{Field: "name", Order: "ASC"}, {Field: "id", Order: "DESC"}}:
//   Устанавливает qb.OrderByClause = " ORDER BY name ASC, id DESC"
// Возвращаемое значение: void
func (qb *QueryBuilder) AddSort(sorts []SortBy, fieldValidator func(string) bool) {
	var orderBy []string
	for _, sort := range sorts {
		if fieldValidator(sort.Field) {
			orderBy = append(orderBy, fmt.Sprintf("%s %s", sort.Field, sort.Order))
		}
	}

	if len(orderBy) > 0 {
		qb.OrderByClause = " ORDER BY " + strings.Join(orderBy, ", ")
	}
}

// AddPagination добавляет пагинацию LIMIT и OFFSET к запросу.
// page: номер страницы (начинается с 1)
// size: количество элементов на странице
// Пример для page=2, size=10:
//   Устанавливает qb.LimitClause = " LIMIT $3"
//   Устанавливает qb.OffsetClause = " OFFSET $4"
//   Добавляет параметры: 10 и 10 (10 = (2-1)*10) в qb.Params
//   Увеличивает paramCounter на 2
// Если page=1, OFFSET не добавляется.
func (qb *QueryBuilder) AddPagination(page, size int) {
	if size > 0 {
		qb.LimitClause = fmt.Sprintf(" LIMIT $%d", qb.paramCounter)
		qb.Params = append(qb.Params, size)
		qb.paramCounter++

		if page > 1 {
			qb.OffsetClause = fmt.Sprintf(" OFFSET $%d", qb.paramCounter)
			qb.Params = append(qb.Params, (page-1)*size)
			qb.paramCounter++
		}
	}
}

// BuildWhereClause собирает все условия WHERE в одну строку.
// Возвращает пустую строку, если условий нет.
// Пример при наличии условий ["name = $1", "id > $2"]:
//   Возвращает: " WHERE name = $1 AND id > $2"
// Пример при отсутствии условий:
//   Возвращает: ""
func (qb *QueryBuilder) BuildWhereClause() string {
	if len(qb.WhereConditions) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(qb.WhereConditions, " AND ")
}

// BuildSelectQuery строит полный SELECT запрос.
// baseQuery: базовая часть SELECT запроса (например: "SELECT * FROM users")
// Пример для baseQuery="SELECT * FROM users" с условиями WHERE и сортировкой:
//   Возвращает: "SELECT * FROM users WHERE name = $1 ORDER BY id DESC LIMIT $2 OFFSET $3"
func (qb *QueryBuilder) BuildSelectQuery(baseQuery string) string {
	query := baseQuery
	query += qb.BuildWhereClause()
	query += qb.OrderByClause
	query += qb.LimitClause
	query += qb.OffsetClause
	return query
}

// BuildCountQuery строит COUNT запрос для подсчета общего количества записей.
// baseQuery: базовая часть SELECT запроса без пагинации
// Пример для baseQuery="SELECT * FROM users" с условиями WHERE:
//   Возвращает: "SELECT COUNT(*) FROM (SELECT * FROM users WHERE name = $1) as subquery"
// не включает LIMIT и OFFSET
func (qb *QueryBuilder) BuildCountQuery(baseQuery string) string {
	query := "SELECT COUNT(*) FROM (" + baseQuery + qb.BuildWhereClause() + ") as subquery"
	return query
}

// GetParams возвращает все параметры запроса включая параметры пагинации.
// Пример после добавления нескольких фильтров и пагинации:
//   Возвращает: []interface{}{"test", 100, 10, 10} // значения для $1, $2, $3, $4
func (qb *QueryBuilder) GetParams() []interface{} {
	return qb.Params
}

// GetCountParams возвращает параметры для COUNT запроса (без параметров пагинации).
// Пример при наличии параметров: []interface{}{"test", 100, 10, 10}
//   Возвращает: []interface{}{"test", 100} // только параметры для WHERE условий
func (qb *QueryBuilder) GetCountParams() []interface{} {
	paramCount := len(qb.Params)
	if qb.LimitClause != "" {
		paramCount--
	}
	if qb.OffsetClause != "" {
		paramCount--
	}

	if paramCount > 0 {
		return qb.Params[:paramCount]
	}
	return []interface{}{}
}

type ListQuery struct {
	SelectQuery  string
	CountQuery   string
	SelectParams []interface{}
	CountParams  []interface{}
}

func BuildListQuery(req PagebleRq, baseSelectQuery string, fieldValidator func(string) bool) *ListQuery {
	qb := NewQueryBuilder()

	for field, predicate := range req.Filter {
		qb.AddFilter(field, predicate, fieldValidator)
	}

	qb.AddSort(req.Sort, fieldValidator)

	qb.AddPagination(req.Page, req.Size)

	return &ListQuery{
		SelectQuery:  qb.BuildSelectQuery(baseSelectQuery),
		CountQuery:   qb.BuildCountQuery(baseSelectQuery),
		SelectParams: qb.GetParams(),
		CountParams:  qb.GetCountParams(),
	}
}
