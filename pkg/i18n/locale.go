package i18n

import (
	"context"
	"net/http"
	"strings"
)

type Locale string

const (
	LocaleRU Locale = "ru"
	LocaleEN Locale = "en"
)

type contextKey struct{}

// WithLocale сохраняет локаль в context.
func WithLocale(ctx context.Context, locale Locale) context.Context {
	return context.WithValue(ctx, contextKey{}, locale)
}

// FromContext возвращает локаль из context или LocaleRU по умолчанию.
func FromContext(ctx context.Context) Locale {
	if locale, ok := ctx.Value(contextKey{}).(Locale); ok && locale.IsValid() {
		return locale
	}
	return LocaleRU
}

func (l Locale) IsValid() bool {
	return l == LocaleRU || l == LocaleEN
}

// ResolveLocale определяет локаль из заголовка Accept-Language.
func ResolveLocale(r *http.Request) Locale {
	header := strings.TrimSpace(r.Header.Get("Accept-Language"))
	if header == "" {
		return LocaleRU
	}
	primary := strings.ToLower(strings.TrimSpace(strings.Split(header, ",")[0]))
	primary = strings.TrimSpace(strings.Split(primary, ";")[0])
	if strings.HasPrefix(primary, "en") {
		return LocaleEN
	}
	return LocaleRU
}
