package entityreqdecorator

import (
	"strconv"
	"strings"
)

const (
	PAGE = 1
	SIZE = 10
)

func ParseQueryParams(queryParams map[string][]string) PagebleRq {
	req := PagebleRq{
		Page:   PAGE,
		Size:   SIZE,
		Sort:   []SortBy{},
		Filter: make(map[string]SQLGenerator),
	}

	// Парсинг пагинации
	if pageStr, ok := queryParams["page"]; ok && len(pageStr) > 0 {
		if page, err := strconv.Atoi(pageStr[0]); err == nil && page > 0 {
			req.Page = page
		}
	}

	if sizeStr, ok := queryParams["size"]; ok && len(sizeStr) > 0 {
		if size, err := strconv.Atoi(sizeStr[0]); err == nil && size > 0 {
			req.Size = size
		}
	}

	if sortParams, ok := queryParams["sort"]; ok {
		for _, sortParam := range sortParams {
			if parts := strings.Split(sortParam, ","); len(parts) == 2 {
				field := strings.TrimSpace(parts[0])
				order := strings.TrimSpace(strings.ToUpper(parts[1]))
				if order == "ASC" || order == "DESC" {
					req.Sort = append(req.Sort, SortBy{
						Field: field,
						Order: order,
					})
				}
			}
		}
	}

	// Парсинг фильтров
	for key, values := range queryParams {
		// Пропускаем уже обработанные параметры
		if key == "page" || key == "size" || key == "sort" {
			continue
		}

		if len(values) > 0 {
			req.Filter[key] = parsePredicate(key, values[0])
		}
	}

	return req
}

func parsePredicate(field, value string) SQLGenerator {
	if value == "" {
		return &PredicateEQ{
			Predicate: Predicate{
				Field: field,
				Value: "",
			},
		}
	}

	if !strings.Contains(value, "(") || !strings.HasSuffix(value, ")") {
		return &PredicateEQ{
			Predicate: Predicate{
				Field: field,
				Value: value,
			},
		}
	}

	operator, innerValue := extractOperator(value)

	switch operator {
	// TODO or predicate
	case "like":
		return &PredicateLike{Predicate: Predicate{Value: innerValue, Field: field}}
	case "anf":
		// anf(gt(20),lt(30)) - AND
		predicates := parseANF(innerValue, field)
		return &PredicateANF{
			Predicate: Predicate{
				Field:          field,
				InnerPredicate: predicates,
			},
		}
	case "gt":
		return &PredicateGT{
			Predicate: Predicate{
				Field: field,
				Value: innerValue,
			},
		}
	case "lt":
		return &PredicateLT{
			Predicate: Predicate{
				Field: field,
				Value: innerValue,
			},
		}
	case "gte":
		return &PredicateGTE{
			Predicate: Predicate{
				Field: field,
				Value: innerValue,
			},
		}
	case "lte":
		return &PredicateLTE{
			Predicate: Predicate{
				Field: field,
				Value: innerValue,
			},
		}
	case "ne":
		return &PredicateNE{
			Predicate: Predicate{
				Field: field,
				Value: innerValue,
			},
		}
	default:
		return &PredicateEQ{
			Predicate: Predicate{
				Field: field,
				Value: value, // сохраняем оригинальное значение
			},
		}
	}
}
func extractOperator(s string) (string, string) {
	openParen := strings.Index(s, "(")
	if openParen == -1 {
		return "", s
	}

	operator := strings.TrimSpace(s[:openParen])
	inner := s[openParen+1 : len(s)-1] // убираем внешние скобки

	return operator, inner
}
func parseANF(s, field string) []SQLGenerator {
	var predicates []SQLGenerator
	parts := splitByCommaOutsideParens(s)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		operator, innerValue := extractOperator(part)

		switch operator {
		case "gt":
			predicates = append(predicates, &PredicateGT{
				Predicate: Predicate{
					Value: innerValue,
					Field: field,
				},
			})
		case "lt":
			predicates = append(predicates, &PredicateLT{
				Predicate: Predicate{
					Value: innerValue, Field: field,
				},
			})
		case "gte":
			predicates = append(predicates, &PredicateGTE{
				Predicate: Predicate{
					Value: innerValue, Field: field,
				},
			})
		case "lte":
			predicates = append(predicates, &PredicateLTE{
				Predicate: Predicate{
					Value: innerValue, Field: field,
				},
			})
		case "eq":
			predicates = append(predicates, &PredicateEQ{
				Predicate: Predicate{
					Value: innerValue, Field: field,
				},
			})
		case "like":
			predicates = append(predicates, &PredicateLike{
				Predicate: Predicate{
					Value: innerValue, Field: field,
				},
			})
		}
	}

	return predicates
}

func splitByCommaOutsideParens(s string) []string {
	var result []string
	var current strings.Builder
	parenDepth := 0

	for _, r := range s {
		switch r {
		case '(':
			parenDepth++
			current.WriteRune(r)
		case ')':
			parenDepth--
			current.WriteRune(r)
		case ',':
			if parenDepth == 0 {
				result = append(result, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
