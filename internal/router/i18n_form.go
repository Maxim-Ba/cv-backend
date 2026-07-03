package router

import (
	"net/http"

	"github.com/Maxim-Ba/cv-backend/internal/models/gen"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
)

func localizedFromForm(r *http.Request, name string) i18n.LocalizedText {
	return i18n.ParseFormLocalized(r, name)
}

func localizedListFromForm(r *http.Request, name string) i18n.LocalizedStringList {
	ru := i18n.ParseLinesInput(r.FormValue(name + "_ru"))
	en := i18n.ParseLinesInput(r.FormValue(name + "_en"))
	return i18n.NewLocalizedStringList(ru, en)
}

func nullableLocalizedFromForm(r *http.Request, name string) i18n.NullableLocalizedText {
	text := i18n.ParseFormLocalized(r, name)
	if text.IsEmpty() {
		return i18n.NullableLocalizedText{}
	}
	return i18n.NullableLocalizedText{Text: text, Valid: true}
}

func tagFromForm(r *http.Request, id int64, hexColor string) models.Tag {
	return models.Tag{
		ID:       id,
		Name:     localizedFromForm(r, "name"),
		HexColor: hexColor,
	}
}

func educationFromForm(r *http.Request, id int64, year int32) models.Education {
	return models.Education{
		ID:           id,
		Name:         nullableLocalizedFromForm(r, "name"),
		Year:         year,
		Course:       localizedFromForm(r, "course"),
		Organization: localizedFromForm(r, "organization"),
	}
}

func workHistoryFromForm(r *http.Request, wh models.WorkHistory) models.WorkHistory {
	wh.Name = localizedFromForm(r, "name")
	wh.JobTitle = nullableLocalizedFromForm(r, "jobTitle")
	wh.About = localizedFromForm(r, "about")
	wh.WhatIDid = localizedListFromForm(r, "whatIDid")
	wh.Projects = localizedListFromForm(r, "projects")
	return wh
}

func technologyDescriptionFromForm(r *http.Request) i18n.NullableLocalizedText {
	return nullableLocalizedFromForm(r, "description")
}
