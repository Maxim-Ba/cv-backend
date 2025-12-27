package entityreqdecorator

import (
	"slices"
	"testing"
)

func TestParseQueryParams(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		queryParams map[string][]string
		want        PagebleRq
	}{
		{
			name: "parameters without operators",
			queryParams: map[string][]string{
				"size": {"11"},
				"page": {"2"},
			},
			want: PagebleRq{
				Size: 11,
				Page: 2,
			},
		},
		{
			name:        "empty parameters",
			queryParams: map[string][]string{},
			want: PagebleRq{
				Size: SIZE,
				Page: PAGE,
			},
		},
		{
			name: "parameters with simple operator",
			queryParams: map[string][]string{
				"name": {"like(vasya)"},
			},
			want: PagebleRq{
				Size: SIZE,
				Page: PAGE,
				Filter: map[string]SQLGenerator{
					"name": &PredicateLike{
						Predicate: Predicate{
							Value: "vasya",
							Field: "name",
						},
					},
				},
			},
		},
		{
			name: "parameters with sorting",
			queryParams: map[string][]string{
				"name": {"like(vasya)"},
				"sort": {"id,ASC"},
			},
			want: PagebleRq{
				Size: SIZE,
				Page: PAGE,
				Sort: []SortBy{
					{
						Field: "id",
						Order: "ASC",
					},
				},
				Filter: map[string]SQLGenerator{
					"name": &PredicateLike{
						Predicate: Predicate{
							Value: "vasya",
							Field: "name",
						},
					},
				},
			},
		},
		{
			name: "parameters with multisorting",
			queryParams: map[string][]string{
				"name": {"like(vasya)"},
				"sort": {"id,ASC", "name,ASC"},
			},
			want: PagebleRq{
				Size: SIZE,
				Page: PAGE,
				Sort: []SortBy{
					{
						Field: "id",
						Order: "ASC",
					},
					{
						Field: "name",
						Order: "ASC",
					},
				},
				Filter: map[string]SQLGenerator{
					"name": &PredicateLike{
						Predicate: Predicate{
							Value: "vasya",
							Field: "name",
						},
					},
				},
			},
		},
		{
			name: "complex parameters",
			queryParams: map[string][]string{
				"date": {"anf(gt(10-10-2025),lt(20-10-2025))"},
			},
			want: PagebleRq{
				Size: SIZE,
				Page: PAGE,
				Filter: map[string]SQLGenerator{
					"date": &PredicateANF{
						Predicate: Predicate{
							InnerPredicate: []SQLGenerator{
								&PredicateGT{
									Predicate: Predicate{
										Value: "10-10-2025",
										Field: "date",
									},
								},
								&PredicateLT{
									Predicate: Predicate{
										Value: "20-10-2025",
										Field: "date",
									},
								},
							},
							Field: "date",
						},
					},
				},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseQueryParams(tt.queryParams)
			// Size and Page check
			isEqualPageSize := got.Page == tt.want.Page && got.Size == tt.want.Size
			if !isEqualPageSize {
				t.Errorf("ParseQueryParams() got.Page= %d, got.Size= %d, want.Page= %d want.Size= %d", got.Page, got.Size, tt.want.Page, tt.want.Size)
			}
			// Sort check
			if len(tt.want.Sort) != 0 {
				if len(tt.want.Sort) != len(got.Sort) {
					t.Errorf("got.SortLength %d, want.Sort length %d", len(got.Sort), len(tt.want.Sort))
				}
				for _, wantSort := range tt.want.Sort {
					has := slices.Contains(got.Sort, wantSort)
					if !has {
						t.Errorf("wantSort %v not in got.Sort %v", wantSort, got.Sort)
					}
				}
			}
			// Filter check
			if len(tt.want.Filter) != 0 {
				if len(tt.want.Filter) != len(got.Filter) {
					t.Errorf("got.FilterLength %d, want.Filter length %d", len(got.Filter), len(tt.want.Filter))
				}
				for key, wantFilter := range tt.want.Filter {
					gotFilter, has := got.Filter[key]
					if !has {
						t.Errorf(
							"wantFilter %v not in got.Filter %v",
							wantFilter,
							got.Filter,
						)
					}
					if gotFilter.GenerateSQL() != wantFilter.GenerateSQL() {
						t.Errorf(
							"gotFilter %v with gotFilter.GenerateSQL() %s not equal wantFilter %v with wantFilter.GenerateSQL() %s",
							gotFilter, gotFilter.GenerateSQL(),
							wantFilter, wantFilter.GenerateSQL(),
						)
					}
				}
			}
		})
	}
}
