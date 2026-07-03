package middleware

import (
	"net/http"

	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
)

// LocaleMiddleware сохраняет локаль из Accept-Language в context запроса.
func LocaleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := i18n.ResolveLocale(r)
		ctx := i18n.WithLocale(r.Context(), locale)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
