package dto

// TagDTO — плоское представление тега для JSON-ответа
type TagDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

// TechnologyDTO — плоское представление технологии для JSON-ответа (без pgtype)
type TechnologyDTO struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	LogoUrl     *string `json:"logoUrl"`
}

// TechnologyWithTagsDTO — технология с вложенными тегами
type TechnologyWithTagsDTO struct {
	TechnologyDTO
	Tags []TagDTO `json:"tags"`
}

// WorkHistoryDTO — плоское представление истории работы (даты как строки ISO, logoUrl как строка)
type WorkHistoryDTO struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	JobTitle    string   `json:"jobTitle"`
	About       string   `json:"about"`
	LogoUrl     string   `json:"logoUrl"`
	PeriodStart *string  `json:"periodStart"`
	PeriodEnd   *string  `json:"periodEnd"`
	WhatIDid    []string `json:"whatIDid"`
	Projects    []string `json:"projects"`
}

// WorkHistoryWithTechnologiesDTO — история работы с вложенными технологиями
type WorkHistoryWithTechnologiesDTO struct {
	WorkHistoryDTO
	Technologies []TechnologyWithTagsDTO `json:"technologies"`
}

// EducationDTO — плоское представление образования для JSON-ответа
type EducationDTO struct {
	ID           int64   `json:"id"`
	Name         *string `json:"name"`
	Year         int32   `json:"year"`
	Course       string  `json:"course"`
	Organization string  `json:"organization"`
}

// TagListResponse — страничный ответ для тегов (swagger)
type TagListResponse struct {
	Total   int      `json:"total,omitempty"`
	Content []TagDTO `json:"content"`
	Page    int      `json:"page,omitempty"`
	Size    int      `json:"size,omitempty"`
}

// TechListResponse — страничный ответ для технологий (swagger)
type TechListResponse struct {
	Total   int                    `json:"total,omitempty"`
	Content []TechnologyWithTagsDTO `json:"content"`
	Page    int                    `json:"page,omitempty"`
	Size    int                    `json:"size,omitempty"`
}

// WorkHistoryListResponse — страничный ответ для истории работы (swagger)
type WorkHistoryListResponse struct {
	Total   int                              `json:"total,omitempty"`
	Content []WorkHistoryWithTechnologiesDTO `json:"content"`
	Page    int                              `json:"page,omitempty"`
	Size    int                              `json:"size,omitempty"`
}

// EducationListResponse — страничный ответ для образования (swagger)
type EducationListResponse struct {
	Total   int           `json:"total,omitempty"`
	Content []EducationDTO `json:"content"`
	Page    int           `json:"page,omitempty"`
	Size    int           `json:"size,omitempty"`
}

// DeleteResponse — ответ при удалении записей (swagger)
type DeleteResponse struct {
	DeletedIDs []int64 `json:"deleted_ids"`
	Count      int     `json:"count"`
}

// AboutMeDTO — данные секции «О себе» для JSON-ответа
type AboutMeDTO struct {
	BioParagraphs []string                `json:"bioParagraphs"`
	Technologies  []TechnologyWithTagsDTO `json:"technologies"`
	Note          *string                 `json:"note"`
	Hobbies       *string                 `json:"hobbies"`
}

// HeroDTO — данные hero-секции для JSON-ответа
type HeroDTO struct {
	Greeting string `json:"greeting"`
	FullName string `json:"fullName"`
	Title    string `json:"title"`
	Pitch    string `json:"pitch"`
}
