package services

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Maxim-Ba/cv-backend/internal/models/dto"
	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
)

// MockWorkHistoryRepo мок-репозиторий для тестирования WorkHistoryService
type MockWorkHistoryRepo struct {
	GetFunc        func(id int64) (models.WorkHistory, error)
	ListFunc       func(entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.WorkHistory], error)
	CreateFunc     func(models.WorkHistory) (models.WorkHistory, error)
	UpdateFunc     func(models.WorkHistory) (models.WorkHistory, error)
	DeleteFunc     func(id int64) (int64, error)
	DeleteListFunc func([]int64) ([]int64, error)
}

func (m *MockWorkHistoryRepo) Get(id int64) (models.WorkHistory, error) {
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return models.WorkHistory{}, nil
}

func (m *MockWorkHistoryRepo) List(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.WorkHistory], error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return entityreqdecorator.PagebleRs[models.WorkHistory]{}, nil
}

func (m *MockWorkHistoryRepo) Create(wh models.WorkHistory) (models.WorkHistory, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(wh)
	}
	return models.WorkHistory{}, nil
}

func (m *MockWorkHistoryRepo) Update(wh models.WorkHistory) (models.WorkHistory, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(wh)
	}
	return models.WorkHistory{}, nil
}

func (m *MockWorkHistoryRepo) Delete(id int64) (int64, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return 0, nil
}

func (m *MockWorkHistoryRepo) DeleteList(ids []int64) ([]int64, error) {
	if m.DeleteListFunc != nil {
		return m.DeleteListFunc(ids)
	}
	return nil, nil
}

func (m *MockWorkHistoryRepo) GetWithTechnologies(id int64, locale i18n.Locale) (dto.WorkHistoryWithTechnologiesDTO, error) {
	return dto.WorkHistoryWithTechnologiesDTO{}, nil
}

func (m *MockWorkHistoryRepo) ListWithTechnologies(req entityreqdecorator.PagebleRq, locale i18n.Locale) (entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO], error) {
	return entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO]{}, nil
}

func (m *MockWorkHistoryRepo) SetTechnologies(workHistoryID int64, technologyIDs []int64) error {
	return nil
}

// TestWorkHistoryService_Get тестирует метод Get
func TestWorkHistoryService_Get(t *testing.T) {
	testDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        int64
		mockWH    models.WorkHistory
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное получение истории работы",
			id:   1,
			mockWH: models.WorkHistory{
				ID:          1,
				Name: newLocalizedText("Яндекс"),
				About: newLocalizedText("Работал Backend разработчиком"),
				LogoUrl: newPgText("logo.png"),
				PeriodStart: pgtype.Date{Time: testDate, Valid: true},
				PeriodEnd:   pgtype.Date{Time: testDate.AddDate(2, 0, 0), Valid: true},
				WhatIDid: newLocalizedStringList("Разработка API", "Оптимизация БД"),
				Projects: newLocalizedStringList("Поиск", "Карты"),
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Невалидный ID (0)",
			id:        0,
			mockError: nil,
			wantError: true,
			errorMsg:  "invalid work history ID",
		},
		{
			name:      "Ошибка репозитория",
			id:        1,
			mockError: errors.New("database error"),
			wantError: true,
			errorMsg:  "error getting work history",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockWorkHistoryRepo{
				GetFunc: func(id int64) (models.WorkHistory, error) {
					return tt.mockWH, tt.mockError
				},
			}
			service := NewWorkHistoryService(mockRepo)

			// Act
			result, err := service.Get(tt.id)

			// Assert
			if tt.wantError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но получили nil")
				}
				if tt.errorMsg != "" && err != nil {
					if !contains(err.Error(), tt.errorMsg) {
						t.Errorf("Ожидалось сообщение об ошибке содержащее '%s', получили: %v", tt.errorMsg, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Не ожидалась ошибка, получили: %v", err)
				}
				if result.ID != tt.mockWH.ID {
					t.Errorf("Ожидался ID = %d, получили %d", tt.mockWH.ID, result.ID)
				}
				if result.Name.Get(i18n.LocaleRU) != tt.mockWH.Name.Get(i18n.LocaleRU) {
					t.Errorf("Ожидалось Name = %s, получили %s", tt.mockWH.Name.Get(i18n.LocaleRU), result.Name.Get(i18n.LocaleRU))
				}
				if len(result.WhatIDid.Get(i18n.LocaleRU)) != len(tt.mockWH.WhatIDid.Get(i18n.LocaleRU)) {
					t.Errorf("Ожидалось %d элементов WhatIDid, получили %d", len(tt.mockWH.WhatIDid.Get(i18n.LocaleRU)), len(result.WhatIDid.Get(i18n.LocaleRU)))
				}
			}
		})
	}
}

// TestWorkHistoryService_List тестирует метод List
func TestWorkHistoryService_List(t *testing.T) {
	testDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		request    entityreqdecorator.PagebleRq
		mockResult entityreqdecorator.PagebleRs[models.WorkHistory]
		mockError  error
		wantError  bool
	}{
		{
			name: "Успешное получение списка историй работы",
			request: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 10,
			},
			mockResult: entityreqdecorator.PagebleRs[models.WorkHistory]{
				Total: 2,
				Content: []models.WorkHistory{
					{
						ID:          1,
						Name: newLocalizedText("Яндекс"),
						About: newLocalizedText("Backend разработчик"),
						PeriodStart: pgtype.Date{Time: testDate, Valid: true},
						PeriodEnd:   pgtype.Date{Time: testDate.AddDate(2, 0, 0), Valid: true},
						WhatIDid: newLocalizedStringList("API", "БД"),
						Projects: newLocalizedStringList("Поиск"),
					},
					{
						ID:          2,
						Name: newLocalizedText("Google"),
						About: newLocalizedText("Software Engineer"),
						PeriodStart: pgtype.Date{Time: testDate.AddDate(-3, 0, 0), Valid: true},
						PeriodEnd:   pgtype.Date{Time: testDate, Valid: true},
						WhatIDid: newLocalizedStringList("Development", "Testing"),
						Projects: newLocalizedStringList("Gmail", "Drive"),
					},
				},
				Page: 1,
				Size: 10,
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Ошибка репозитория",
			request: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 10,
			},
			mockError: errors.New("connection timeout"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockWorkHistoryRepo{
				ListFunc: func(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.WorkHistory], error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewWorkHistoryService(mockRepo)

			// Act
			result, err := service.List(tt.request)

			// Assert
			if tt.wantError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но получили nil")
				}
			} else {
				if err != nil {
					t.Errorf("Не ожидалась ошибка, получили: %v", err)
				}
				if result.Total != tt.mockResult.Total {
					t.Errorf("Ожидался Total = %d, получили %d", tt.mockResult.Total, result.Total)
				}
				if len(result.Content) != len(tt.mockResult.Content) {
					t.Errorf("Ожидалось %d элементов, получили %d", len(tt.mockResult.Content), len(result.Content))
				}
			}
		})
	}
}

// TestWorkHistoryService_Create тестирует метод Create
func TestWorkHistoryService_Create(t *testing.T) {
	testDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		wh        models.WorkHistory
		mockWH    models.WorkHistory
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное создание истории работы",
			wh: models.WorkHistory{
				Name: newLocalizedText("Тинькофф"),
				About: newLocalizedText("Go разработчик"),
				LogoUrl: newPgText("tinkoff.png"),
				PeriodStart: pgtype.Date{Time: testDate, Valid: true},
				PeriodEnd:   pgtype.Date{Time: testDate.AddDate(1, 0, 0), Valid: true},
				WhatIDid: newLocalizedStringList("Микросервисы", "Kafka"),
				Projects: newLocalizedStringList("Банкинг", "Инвестиции"),
			},
			mockWH: models.WorkHistory{
				ID:          3,
				Name: newLocalizedText("Тинькофф"),
				About: newLocalizedText("Go разработчик"),
				LogoUrl: newPgText("tinkoff.png"),
				PeriodStart: pgtype.Date{Time: testDate, Valid: true},
				PeriodEnd:   pgtype.Date{Time: testDate.AddDate(1, 0, 0), Valid: true},
				WhatIDid: newLocalizedStringList("Микросервисы", "Kafka"),
				Projects: newLocalizedStringList("Банкинг", "Инвестиции"),
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Отсутствует имя компании",
			wh: models.WorkHistory{
				Name:  newLocalizedText(""),
				About: newLocalizedText("Some description"),
			},
			wantError: true,
			errorMsg:  "name and about are required fields",
		},
		{
			name: "Отсутствует описание",
			wh: models.WorkHistory{
				Name: newLocalizedText("Company"),
				About: newLocalizedText(""),
			},
			wantError: true,
			errorMsg:  "name and about are required fields",
		},
		{
			name: "Ошибка репозитория",
			wh: models.WorkHistory{
				Name: newLocalizedText("Company"),
				About: newLocalizedText("Description"),
			},
			mockError: errors.New("database constraint violation"),
			wantError: true,
			errorMsg:  "error creating work history",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockWorkHistoryRepo{
				CreateFunc: func(wh models.WorkHistory) (models.WorkHistory, error) {
					return tt.mockWH, tt.mockError
				},
			}
			service := NewWorkHistoryService(mockRepo)

			// Act
			result, err := service.Create(tt.wh)

			// Assert
			if tt.wantError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но получили nil")
				}
				if tt.errorMsg != "" && err != nil {
					if !contains(err.Error(), tt.errorMsg) {
						t.Errorf("Ожидалось сообщение об ошибке содержащее '%s', получили: %v", tt.errorMsg, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Не ожидалась ошибка, получили: %v", err)
				}
				if result.ID != tt.mockWH.ID {
					t.Errorf("Ожидался ID = %d, получили %d", tt.mockWH.ID, result.ID)
				}
				if result.Name.Get(i18n.LocaleRU) != tt.mockWH.Name.Get(i18n.LocaleRU) {
					t.Errorf("Ожидалось Name = %s, получили %s", tt.mockWH.Name.Get(i18n.LocaleRU), result.Name.Get(i18n.LocaleRU))
				}
			}
		})
	}
}

// TestWorkHistoryService_Update тестирует метод Update
func TestWorkHistoryService_Update(t *testing.T) {
	testDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		wh        models.WorkHistory
		mockWH    models.WorkHistory
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное обновление истории работы",
			wh: models.WorkHistory{
				ID:          1,
				Name: newLocalizedText("Яндекс"),
				About: newLocalizedText("Senior Backend Developer"),
				PeriodStart: pgtype.Date{Time: testDate, Valid: true},
				PeriodEnd:   pgtype.Date{Time: testDate.AddDate(3, 0, 0), Valid: true},
				WhatIDid: newLocalizedStringList("Lead development", "Mentoring"),
				Projects: newLocalizedStringList("Search", "Maps", "Cloud"),
			},
			mockWH: models.WorkHistory{
				ID:          1,
				Name: newLocalizedText("Яндекс"),
				About: newLocalizedText("Senior Backend Developer"),
				PeriodStart: pgtype.Date{Time: testDate, Valid: true},
				PeriodEnd:   pgtype.Date{Time: testDate.AddDate(3, 0, 0), Valid: true},
				WhatIDid: newLocalizedStringList("Lead development", "Mentoring"),
				Projects: newLocalizedStringList("Search", "Maps", "Cloud"),
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Невалидный ID (0)",
			wh: models.WorkHistory{
				ID:    0,
				Name: newLocalizedText("Company"),
				About: newLocalizedText("Description"),
			},
			wantError: true,
			errorMsg:  "invalid work history ID",
		},
		{
			name: "Отсутствует имя",
			wh: models.WorkHistory{
				ID:    1,
				Name:  newLocalizedText(""),
				About: newLocalizedText("Description"),
			},
			wantError: true,
			errorMsg:  "name and about are required fields",
		},
		{
			name: "Ошибка репозитория - история не найдена",
			wh: models.WorkHistory{
				ID:    999,
				Name: newLocalizedText("Company"),
				About: newLocalizedText("Description"),
			},
			mockError: errors.New("work history not found"),
			wantError: true,
			errorMsg:  "error updating work history",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockWorkHistoryRepo{
				UpdateFunc: func(wh models.WorkHistory) (models.WorkHistory, error) {
					return tt.mockWH, tt.mockError
				},
			}
			service := NewWorkHistoryService(mockRepo)

			// Act
			result, err := service.Update(tt.wh)

			// Assert
			if tt.wantError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но получили nil")
				}
				if tt.errorMsg != "" && err != nil {
					if !contains(err.Error(), tt.errorMsg) {
						t.Errorf("Ожидалось сообщение об ошибке содержащее '%s', получили: %v", tt.errorMsg, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Не ожидалась ошибка, получили: %v", err)
				}
				if result.About.Get(i18n.LocaleRU) != tt.mockWH.About.Get(i18n.LocaleRU) {
					t.Errorf("Ожидался About = %s, получили %s", tt.mockWH.About.Get(i18n.LocaleRU), result.About.Get(i18n.LocaleRU))
				}
				if len(result.Projects) != len(tt.mockWH.Projects) {
					t.Errorf("Ожидалось %d проектов, получили %d", len(tt.mockWH.Projects), len(result.Projects))
				}
			}
		})
	}
}

// TestWorkHistoryService_Delete тестирует метод Delete
func TestWorkHistoryService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         int64
		mockResult int64
		mockError  error
		wantError  bool
		errorMsg   string
	}{
		{
			name:       "Успешное удаление",
			id:         1,
			mockResult: 1,
			mockError:  nil,
			wantError:  false,
		},
		{
			name:      "Невалидный ID (0)",
			id:        0,
			wantError: true,
			errorMsg:  "invalid work history ID",
		},
		{
			name:      "Ошибка репозитория",
			id:        1,
			mockError: errors.New("foreign key constraint"),
			wantError: true,
			errorMsg:  "error deleting work history",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockWorkHistoryRepo{
				DeleteFunc: func(id int64) (int64, error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewWorkHistoryService(mockRepo)

			// Act
			result, err := service.Delete(tt.id)

			// Assert
			if tt.wantError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но получили nil")
				}
				if tt.errorMsg != "" && err != nil {
					if !contains(err.Error(), tt.errorMsg) {
						t.Errorf("Ожидалось сообщение об ошибке содержащее '%s', получили: %v", tt.errorMsg, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Не ожидалась ошибка, получили: %v", err)
				}
				if result != tt.mockResult {
					t.Errorf("Ожидался результат = %d, получили %d", tt.mockResult, result)
				}
			}
		})
	}
}

// TestWorkHistoryService_DeleteList тестирует метод DeleteList
func TestWorkHistoryService_DeleteList(t *testing.T) {
	tests := []struct {
		name       string
		ids        []int64
		mockResult []int64
		mockError  error
		wantError  bool
	}{
		{
			name:       "Успешное удаление списка",
			ids:        []int64{1, 2, 3},
			mockResult: []int64{1, 2, 3},
			mockError:  nil,
			wantError:  false,
		},
		{
			name:       "Пустой список",
			ids:        []int64{},
			mockResult: nil,
			mockError:  nil,
			wantError:  false,
		},
		{
			name:      "Ошибка репозитория",
			ids:       []int64{1, 2},
			mockError: errors.New("database error"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockWorkHistoryRepo{
				DeleteListFunc: func(ids []int64) ([]int64, error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewWorkHistoryService(mockRepo)

			// Act
			result, err := service.DeleteList(tt.ids)

			// Assert
			if tt.wantError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но получили nil")
				}
			} else {
				if err != nil {
					t.Errorf("Не ожидалась ошибка, получили: %v", err)
				}
				if len(result) != len(tt.mockResult) {
					t.Errorf("Ожидалось %d удаленных элементов, получили %d", len(tt.mockResult), len(result))
				}
			}
		})
	}
}
