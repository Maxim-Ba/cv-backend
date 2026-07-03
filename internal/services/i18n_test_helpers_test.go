package services

import (
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
	"github.com/jackc/pgx/v5/pgtype"
)

func newLocalizedText(ru string) i18n.LocalizedText {
	return i18n.FromLegacyText(ru)
}

func newNullableLocalizedText(ru string) i18n.NullableLocalizedText {
	if ru == "" {
		return i18n.NullableLocalizedText{}
	}
	return i18n.NullableLocalizedText{
		Text:  i18n.FromLegacyText(ru),
		Valid: true,
	}
}

func newLocalizedStringList(values ...string) i18n.LocalizedStringList {
	return i18n.FromLegacyStringList(values)
}

func newPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}
