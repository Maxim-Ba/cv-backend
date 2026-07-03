package mapper

import (
	"time"

	"github.com/Maxim-Ba/cv-backend/internal/models/dto"
	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
	"github.com/jackc/pgx/v5/pgtype"
)

func TagToDTO(tag models.Tag, locale i18n.Locale) dto.TagDTO {
	return dto.TagDTO{
		ID:       tag.ID,
		Name:     tag.Name.Resolve(locale),
		HexColor: tag.HexColor,
	}
}

func EducationToDTO(edu models.Education, locale i18n.Locale) dto.EducationDTO {
	return dto.EducationDTO{
		ID:           edu.ID,
		Name:         edu.Name.Resolve(locale),
		Year:         edu.Year,
		Course:       edu.Course.Resolve(locale),
		Organization: edu.Organization.Resolve(locale),
	}
}

func WorkHistoryToDTO(wh models.WorkHistory, locale i18n.Locale) dto.WorkHistoryDTO {
	jobTitle := ""
	if wh.JobTitle.Valid {
		jobTitle = wh.JobTitle.Text.Resolve(locale)
	}
	logoURL := ""
	if wh.LogoUrl.Valid {
		logoURL = wh.LogoUrl.String
	}
	return dto.WorkHistoryDTO{
		ID:          wh.ID,
		Name:        wh.Name.Resolve(locale),
		JobTitle:    jobTitle,
		About:       wh.About.Resolve(locale),
		LogoUrl:     logoURL,
		PeriodStart: dateToString(wh.PeriodStart),
		PeriodEnd:   dateToString(wh.PeriodEnd),
		WhatIDid:    wh.WhatIDid.Resolve(locale),
		Projects:    wh.Projects.Resolve(locale),
	}
}

func TechnologyToDTO(id int64, title string, description i18n.NullableLocalizedText, logoURL pgtype.Text, locale i18n.Locale) dto.TechnologyDTO {
	tech := dto.TechnologyDTO{
		ID:    id,
		Title: title,
	}
	if description.Valid {
		tech.Description = description.Resolve(locale)
	}
	if logoURL.Valid {
		tech.LogoUrl = &logoURL.String
	}
	return tech
}

func dateToString(date pgtype.Date) *string {
	if !date.Valid {
		return nil
	}
	s := date.Time.Format(time.DateOnly)
	return &s
}
