package i18n

import (
	"net/http"
	"strings"
)

// ParseFormLocalized читает пару полей `<name>_ru` и `<name>_en` из HTML-формы.
func ParseFormLocalized(r *http.Request, name string) LocalizedText {
	return NewLocalizedText(
		r.FormValue(name+"_ru"),
		r.FormValue(name+"_en"),
	)
}

// ParseFormLocalizedList читает списки `<name>_ru` и `<name>_en` (по строке на элемент).
func ParseFormLocalizedList(r *http.Request, name string) LocalizedStringList {
	return NewLocalizedStringList(
		r.Form[name+"_ru"],
		r.Form[name+"_en"],
	)
}

// ParseJSONLocalized парсит объект {"ru":"...","en":"..."} из JSON API.
func ParseJSONLocalized(raw map[string]string) LocalizedText {
	if raw == nil {
		return LocalizedText{}
	}
	return LocalizedText{
		string(LocaleRU): strings.TrimSpace(raw[string(LocaleRU)]),
		string(LocaleEN): strings.TrimSpace(raw[string(LocaleEN)]),
	}
}

// ParseJSONLocalizedList парсит объект {"ru":[],"en":[]} из JSON API.
func ParseJSONLocalizedList(raw map[string][]string) LocalizedStringList {
	if raw == nil {
		return LocalizedStringList{}
	}
	return LocalizedStringList{
		string(LocaleRU): normalizeLines(raw[string(LocaleRU)]),
		string(LocaleEN): normalizeLines(raw[string(LocaleEN)]),
	}
}
