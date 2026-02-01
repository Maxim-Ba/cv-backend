package services

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

// MockTechRepo мок-репозиторий для тестирования TechService
type MockTechRepo struct {
	GetFunc        func(id int64) (models.Technology, error)
	ListFunc       func(entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Technology], error)
	CreateFunc     func(models.Technology) (models.Technology, error)
	UpdateFunc     func(models.Technology) (models.Technology, error)
	DeleteFunc     func(id int64) (int64, error)
	DeleteListFunc func([]int64) ([]int64, error)
}

func (m *MockTechRepo) Get(id int64) (models.Technology, error) {
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return models.Technology{}, nil
}

func (m *MockTechRepo) List(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Technology], error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return entityreqdecorator.PagebleRs[models.Technology]{}, nil
}

func (m *MockTechRepo) Create(tech models.Technology) (models.Technology, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(tech)
	}
	return models.Technology{}, nil
}

func (m *MockTechRepo) Update(tech models.Technology) (models.Technology, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(tech)
	}
	return models.Technology{}, nil
}

func (m *MockTechRepo) Delete(id int64) (int64, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return 0, nil
}

func (m *MockTechRepo) DeleteList(ids []int64) ([]int64, error) {
	if m.DeleteListFunc != nil {
		return m.DeleteListFunc(ids)
	}
	return nil, nil
}

// TestTechService_Get тестирует метод Get
func TestTechService_Get(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		mockTech  models.Technology
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное получение технологии",
			id:   1,
			mockTech: models.Technology{
				ID:          1,
				Title:       "Go",
				Description: pgtype.Text{String: "Programming language", Valid: true},
				LogoUrl:     pgtype.Text{String: "https://golang.org/logo.png", Valid: true},
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Невалидный ID (0)",
			id:        0,
			mockError: nil,
			wantError: true,
			errorMsg:  "invalid technology ID",
		},
		{
			name:      "Ошибка репозитория",
			id:        1,
			mockError: errors.New("database error"),
			wantError: true,
			errorMsg:  "error getting technology",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTechRepo{
				GetFunc: func(id int64) (models.Technology, error) {
					return tt.mockTech, tt.mockError
				},
			}
			service := NewTechService(mockRepo)

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
				if result.ID != tt.mockTech.ID {
					t.Errorf("Ожидался ID = %d, получили %d", tt.mockTech.ID, result.ID)
				}
				if result.Title != tt.mockTech.Title {
					t.Errorf("Ожидался Title = %s, получили %s", tt.mockTech.Title, result.Title)
				}
			}
		})
	}
}

// TestTechService_List тестирует метод List
func TestTechService_List(t *testing.T) {
	tests := []struct {
		name       string
		request    entityreqdecorator.PagebleRq
		mockResult entityreqdecorator.PagebleRs[models.Technology]
		mockError  error
		wantError  bool
	}{
		{
			name: "Успешное получение списка технологий",
			request: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 10,
			},
			mockResult: entityreqdecorator.PagebleRs[models.Technology]{
				Total: 2,
				Content: []models.Technology{
					{
						ID:          1,
						Title:       "Go",
						Description: pgtype.Text{String: "Programming language", Valid: true},
					},
					{
						ID:          2,
						Title:       "PostgreSQL",
						Description: pgtype.Text{String: "Database", Valid: true},
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
			mockRepo := &MockTechRepo{
				ListFunc: func(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Technology], error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewTechService(mockRepo)

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

// TestTechService_Create тестирует метод Create
func TestTechService_Create(t *testing.T) {
	tests := []struct {
		name      string
		tech      models.Technology
		mockTech  models.Technology
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное создание технологии",
			tech: models.Technology{
				Title:       "Docker",
				Description: pgtype.Text{String: "Containerization platform", Valid: true},
				LogoUrl:     pgtype.Text{String: "https://docker.com/logo.png", Valid: true},
			},
			mockTech: models.Technology{
				ID:          3,
				Title:       "Docker",
				Description: pgtype.Text{String: "Containerization platform", Valid: true},
				LogoUrl:     pgtype.Text{String: "https://docker.com/logo.png", Valid: true},
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Отсутствует заголовок",
			tech: models.Technology{
				Title:       "",
				Description: pgtype.Text{String: "Some description", Valid: true},
			},
			wantError: true,
			errorMsg:  "technology title is required",
		},
		{
			name: "Ошибка репозитория",
			tech: models.Technology{
				Title: "Docker",
			},
			mockError: errors.New("duplicate key"),
			wantError: true,
			errorMsg:  "error creating technology",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTechRepo{
				CreateFunc: func(tech models.Technology) (models.Technology, error) {
					return tt.mockTech, tt.mockError
				},
			}
			service := NewTechService(mockRepo)

			// Act
			result, err := service.Create(tt.tech)

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
				if result.ID != tt.mockTech.ID {
					t.Errorf("Ожидался ID = %d, получили %d", tt.mockTech.ID, result.ID)
				}
			}
		})
	}
}

// TestTechService_Update тестирует метод Update
func TestTechService_Update(t *testing.T) {
	tests := []struct {
		name      string
		tech      models.Technology
		mockTech  models.Technology
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное обновление технологии",
			tech: models.Technology{
				ID:          1,
				Title:       "Go Updated",
				Description: pgtype.Text{String: "Updated description", Valid: true},
			},
			mockTech: models.Technology{
				ID:          1,
				Title:       "Go Updated",
				Description: pgtype.Text{String: "Updated description", Valid: true},
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Невалидный ID (0)",
			tech: models.Technology{
				ID:    0,
				Title: "Go",
			},
			wantError: true,
			errorMsg:  "invalid technology ID",
		},
		{
			name: "Отсутствует заголовок",
			tech: models.Technology{
				ID:    1,
				Title: "",
			},
			wantError: true,
			errorMsg:  "technology title is required",
		},
		{
			name: "Ошибка репозитория - технология не найдена",
			tech: models.Technology{
				ID:    999,
				Title: "Go",
			},
			mockError: errors.New("technology not found"),
			wantError: true,
			errorMsg:  "error updating technology",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTechRepo{
				UpdateFunc: func(tech models.Technology) (models.Technology, error) {
					return tt.mockTech, tt.mockError
				},
			}
			service := NewTechService(mockRepo)

			// Act
			result, err := service.Update(tt.tech)

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
				if result.Title != tt.mockTech.Title {
					t.Errorf("Ожидался Title = %s, получили %s", tt.mockTech.Title, result.Title)
				}
			}
		})
	}
}

// TestTechService_Delete тестирует метод Delete
func TestTechService_Delete(t *testing.T) {
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
			errorMsg:  "invalid technology ID",
		},
		{
			name:      "Ошибка репозитория",
			id:        1,
			mockError: errors.New("foreign key constraint"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTechRepo{
				DeleteFunc: func(id int64) (int64, error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewTechService(mockRepo)

			// Act
			result, err := service.Delete(tt.id)

			// Assert
			if tt.wantError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но получили nil")
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

// TestTechService_DeleteList тестирует метод DeleteList
func TestTechService_DeleteList(t *testing.T) {
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
			mockRepo := &MockTechRepo{
				DeleteListFunc: func(ids []int64) ([]int64, error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewTechService(mockRepo)

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
