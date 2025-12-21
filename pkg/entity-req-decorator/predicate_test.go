package entityreqdecorator_test

import (
	"testing"

	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

func TestParseQueryParams(t *testing.T) {

	type Case struct {
		name string // description of this test case
		// Named input parameters for target function.
		queryParams map[string][]string
		want        entityreqdecorator.PagebleRq
	}
	tests := []Case{
		Case{
			name: "parameters without operators",
			queryParams: map[string][]string{
				"limit": {"11"},
				"page":  {"2"},
			},
			want: entityreqdecorator.PagebleRq{
				Size: 11,
				Page: 2,
			},
		},
		Case{
			name:        "empty parameters",
			queryParams: map[string][]string{},
			want: entityreqdecorator.PagebleRq{
				Size: entityreqdecorator.SIZE,
				Page: entityreqdecorator.PAGE,
			},
		},
		Case{
			name: "parameters with simple operator",
			queryParams: map[string][]string{
				"name": {"like(vasya)"},
			},
			want: entityreqdecorator.PagebleRq{
				Size: entityreqdecorator.SIZE,
				Page: entityreqdecorator.PAGE,
				Filter: map[string]entityreqdecorator.SQLGenerator{
					"like": &entityreqdecorator.PredicateLike{
						Predicate: entityreqdecorator.Predicate{
							Value: "vasya",
							Field: "name",
						},
					},
				},
			},
		},
		Case{
			name: "parameters with sorting",
			queryParams: map[string][]string{
				"name": {"like(vasya)"},
				"sort": {"id,ASC"},
			},
			want: entityreqdecorator.PagebleRq{
				Size: entityreqdecorator.SIZE,
				Page: entityreqdecorator.PAGE,
				Sort: []entityreqdecorator.SortBy{
					entityreqdecorator.SortBy{
						Field: "id",
						Order: "ASC",
					},
				},
				Filter: map[string]entityreqdecorator.SQLGenerator{
					"like": &entityreqdecorator.PredicateLike{
						Predicate: entityreqdecorator.Predicate{
							Value: "vasya",
							Field: "name",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entityreqdecorator.ParseQueryParams(tt.queryParams)
			if true {
				t.Errorf("ParseQueryParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
