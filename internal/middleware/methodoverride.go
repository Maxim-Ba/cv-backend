package middleware

import "net/http"

// MethodOverride читает скрытое поле _method из form-data
// и подменяет r.Method на PUT или DELETE если указано.
func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err == nil {
				if override := r.FormValue("_method"); override == "PUT" || override == "DELETE" {
					r.Method = override
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
