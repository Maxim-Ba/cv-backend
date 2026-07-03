package i18n

import (
	"net/http"
	"testing"
)

func TestLocalizedTextResolveFallback(t *testing.T) {
	lt := NewLocalizedText("Привет", "")

	if got := lt.Resolve(LocaleRU); got != "Привет" {
		t.Fatalf("ru = %q", got)
	}
	if got := lt.Resolve(LocaleEN); got != "Привет" {
		t.Fatalf("en fallback = %q", got)
	}
}

func TestLocalizedStringListResolveFallback(t *testing.T) {
	ll := NewLocalizedStringList([]string{"один"}, nil)

	if got := ll.Resolve(LocaleEN); len(got) != 1 || got[0] != "один" {
		t.Fatalf("fallback list = %#v", got)
	}
}

func TestResolveLocale(t *testing.T) {
	req := httptestRequest("Accept-Language", "en-US,en;q=0.9")
	if ResolveLocale(req) != LocaleEN {
		t.Fatal("expected en")
	}
}

func httptestRequest(header, value string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set(header, value)
	return req
}
