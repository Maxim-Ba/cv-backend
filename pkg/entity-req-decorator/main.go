package entityreqdecorator

import "strings"

type SortBy struct {
	Field string `json:"field"`
	Order string `json:"order"` // ASC, DESC
}

type PagebleRq struct {
	Page   int 
	Size   int
	Sort   []SortBy 
	Filter map[string]SQLGenerator
}

type PagebleRs[T any] struct {
	Total   int      `json:"total,omitempty"`
	Content []T      `json:"content"`
	Page    int      `json:"page,omitempty"`
	Size    int      `json:"size,omitempty"`
	Sort    []SortBy `json:"sort,omitempty"`
}

type SQLGenerator interface {
	GenerateSQL() string
}
type FieldSetter interface {
	SetField(field string)
}

func init() {
	predicatByStruct = map[string]func(raw string) SQLGenerator{
		"": func(raw string) SQLGenerator {
			return &PredicateEQ{Predicate: Predicate{Value: raw}}
		},
		"eq": func(raw string) SQLGenerator {
			return &PredicateEQ{Predicate: Predicate{Value: raw}}
		},
		"like": func(raw string) SQLGenerator {
			return &PredicateLike{Predicate: Predicate{Value: raw}}
		},
		"gt": func(raw string) SQLGenerator {
			return &PredicateGT{Predicate: Predicate{Value: raw}}
		},
		"lt": func(raw string) SQLGenerator {
			return &PredicateLT{Predicate: Predicate{Value: raw}}
		},
		"gte": func(raw string) SQLGenerator {
			return &PredicateGTE{Predicate: Predicate{Value: raw}}
		},
		"lte": func(raw string) SQLGenerator {
			return &PredicateLTE{Predicate: Predicate{Value: raw}}
		},
		"ne": func(raw string) SQLGenerator {
			return &PredicateNE{Predicate: Predicate{Value: raw}}
		},
		// "anf": func(raw string) SQLGenerator {
		// 	return &PredicateANF{
		// 		predicate: predicate{
		// 			InnerPredicate: parseANF(raw),
		// 		},
		// 	}
		// },
	}
}

var predicatByStruct = map[string]func(raw string) SQLGenerator{
	"": func(raw string) SQLGenerator {
		if raw == "" {
			return &PredicateEQ{}
		}
		return &PredicateEQ{
			Predicate: Predicate{
				Value: raw,
			},
		}
	},
}

type Predicate struct {
	Field          string
	Value          string
	InnerPredicate []SQLGenerator
}
type PredicateLike struct {
	Predicate
}

func (p *PredicateLike) GenerateSQL() string {
	if p.Value == "" {
		return ""
	}
	return p.Field + " LIKE " + p.Value
}

type PredicateEQ struct {
	Predicate
}

func (p *PredicateEQ) GenerateSQL() string {
	if p.Value == "" {
		return ""
	}
	return p.Field + " = " + p.Value
}

type PredicateGT struct {
	Predicate
}

func (p *PredicateGT) GenerateSQL() string {
	if p.Value == "" {
		return ""
	}
	return p.Field + " > " + p.Value
}

func (p *PredicateGT) SetField(field string) {
	p.Field = field
}

type PredicateLT struct {
	Predicate
}

func (p *PredicateLT) GenerateSQL() string {
	if p.Value == "" {
		return ""
	}
	return p.Field + " < " + p.Value
}

func (p *PredicateLT) SetField(field string) {
	p.Field = field
}

type PredicateGTE struct {
	Predicate
}

func (p *PredicateGTE) GenerateSQL() string {
	if p.Value == "" {
		return ""
	}
	return p.Field + " >= " + p.Value
}

func (p *PredicateGTE) SetField(field string) {
	p.Field = field
}

type PredicateLTE struct {
	Predicate
}

func (p *PredicateLTE) GenerateSQL() string {
	if p.Value == "" {
		return ""
	}
	return p.Field + " <= " + p.Value
}

func (p *PredicateLTE) SetField(field string) {
	p.Field = field
}

type PredicateNE struct {
	Predicate
}

func (p *PredicateNE) GenerateSQL() string {
	if p.Value == "" {
		return ""
	}
	return p.Field + " != " + p.Value
}

func (p *PredicateNE) SetField(field string) {
	p.Field = field
}

type PredicateANF struct {
	Predicate
}

func (p *PredicateANF) GenerateSQL() string {
	if len(p.InnerPredicate) == 0 {
		return ""
	}

	var conditions []string
	for _, pred := range p.InnerPredicate {
		if sql := pred.GenerateSQL(); sql != "" {
			conditions = append(conditions, sql)
		}
	}

	if len(conditions) == 0 {
		return ""
	}

	return "(" + strings.Join(conditions, " AND ") + ")"
}

func (p *PredicateANF) SetField(field string) {
	p.Field = field
}
