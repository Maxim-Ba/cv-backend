package components

import (
	"strings"

	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
)

func linesToText(lines []string) string {
	return strings.Join(lines, "\n")
}

func DisplayRU(text i18n.LocalizedText) string {
	return text.Get(i18n.LocaleRU)
}

func DisplayNullableRU(text i18n.NullableLocalizedText) string {
	if !text.Valid {
		return ""
	}
	return text.Text.Get(i18n.LocaleRU)
}

func localizedTextField(name, label string, text i18n.LocalizedText, placeholder string, required, isTextArea bool) FormField {
	return FormField{
		Name:        name,
		Label:       label,
		Type:        "bilingual-textarea",
		ValueRU:     text.Get(i18n.LocaleRU),
		ValueEN:     text.Get(i18n.LocaleEN),
		Placeholder: placeholder,
		Required:    required,
		IsTextArea:  isTextArea,
	}
}

func nullableLocalizedField(name, label string, text i18n.NullableLocalizedText, placeholder string, isTextArea bool) FormField {
	valueRU := ""
	valueEN := ""
	if text.Valid {
		valueRU = text.Text.Get(i18n.LocaleRU)
		valueEN = text.Text.Get(i18n.LocaleEN)
	}
	fieldType := "bilingual-text"
	if isTextArea {
		fieldType = "bilingual-textarea"
	}
	return FormField{
		Name:        name,
		Label:       label,
		Type:        fieldType,
		ValueRU:     valueRU,
		ValueEN:     valueEN,
		Placeholder: placeholder,
		IsTextArea:  isTextArea,
	}
}

func localizedListField(name, label string, list i18n.LocalizedStringList, placeholder string) FormField {
	return FormField{
		Name:        name,
		Label:       label,
		Type:        "bilingual-textarea",
		ValueRU:     linesToText(list.Get(i18n.LocaleRU)),
		ValueEN:     linesToText(list.Get(i18n.LocaleEN)),
		Placeholder: placeholder,
		IsTextArea:  true,
	}
}
