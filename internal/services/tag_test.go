package services

import (
	"errors"
	"testing"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

// MockTagRepo мок-репозиторий для тестирования TagService
type MockTagRepo struct {
	GetFunc        func(id int64) (models.Tag, error)
	ListFunc       func(entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Tag], error)
	CreateFunc     func(models.Tag) (models.Tag, error)
	UpdateFunc     func(models.Tag) (models.Tag, error)
	DeleteFunc     func(id int64) (int64, error)
	DeleteListFunc func([]int64) ([]int64, error)
}

func (m *MockTagRepo) Get(id int64) (models.Tag, error) {
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return models.Tag{}, nil
}

func (m *MockTagRepo) List(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Tag], error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return entityreqdecorator.PagebleRs[models.Tag]{}, nil
}

func (m *MockTagRepo) Create(tag models.Tag) (models.Tag, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(tag)
	}
	return models.Tag{}, nil
}

func (m *MockTagRepo) Update(tag models.Tag) (models.Tag, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(tag)
	}
	return models.Tag{}, nil
}

func (m *MockTagRepo) Delete(id int64) (int64, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return 0, nil
}

func (m *MockTagRepo) DeleteList(ids []int64) ([]int64, error) {
	if m.DeleteListFunc != nil {
		return m.DeleteListFunc(ids)
	}
	return nil, nil
}

// TestTagService_Get тестирует метод Get
func TestTagService_Get(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		mockTag   models.Tag
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное получение тега",
			id:   1,
			mockTag: models.Tag{
				ID:       1,
				Name:     "Backend",
				HexColor: "#FF5733",
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Невалидный ID (0)",
			id:        0,
			mockError: nil,
			wantError: true,
			errorMsg:  "invalid tag ID",
		},
		{
			name:      "Ошибка репозитория",
			id:        1,
			mockError: errors.New("database error"),
			wantError: true,
			errorMsg:  "error getting tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTagRepo{
				GetFunc: func(id int64) (models.Tag, error) {
					return tt.mockTag, tt.mockError
				},
			}
			service := NewTagServise(mockRepo)

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
				if result.ID != tt.mockTag.ID {
					t.Errorf("Ожидался ID = %d, получили %d", tt.mockTag.ID, result.ID)
				}
				if result.Name != tt.mockTag.Name {
					t.Errorf("Ожидалось Name = %s, получили %s", tt.mockTag.Name, result.Name)
				}
			}
		})
	}
}

// TestTagService_List тестирует метод List
func TestTagService_List(t *testing.T) {
	tests := []struct {
		name       string
		request    entityreqdecorator.PagebleRq
		mockResult entityreqdecorator.PagebleRs[models.Tag]
		mockError  error
		wantError  bool
	}{
		{
			name: "Успешное получение списка",
			request: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 10,
			},
			mockResult: entityreqdecorator.PagebleRs[models.Tag]{
				Total: 2,
				Content: []models.Tag{
					{ID: 1, Name: "Backend", HexColor: "#FF5733"},
					{ID: 2, Name: "Frontend", HexColor: "#33FF57"},
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
			mockError: errors.New("database connection error"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTagRepo{
				ListFunc: func(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Tag], error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewTagServise(mockRepo)

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

// TestTagService_Create тестирует метод Create
func TestTagService_Create(t *testing.T) {
	tests := []struct {
		name      string
		tag       models.Tag
		mockTag   models.Tag
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное создание тега",
			tag: models.Tag{
				Name:     "DevOps",
				HexColor: "#00FF00",
			},
			mockTag: models.Tag{
				ID:       3,
				Name:     "DevOps",
				HexColor: "#00FF00",
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Отсутствует имя тега",
			tag: models.Tag{
				Name:     "",
				HexColor: "#00FF00",
			},
			wantError: true,
			errorMsg:  "tag name is required",
		},
		{
			name: "Отсутствует цвет тега",
			tag: models.Tag{
				Name:     "DevOps",
				HexColor: "",
			},
			wantError: true,
			errorMsg:  "tag hex color is required",
		},
		{
			name: "Ошибка репозитория",
			tag: models.Tag{
				Name:     "DevOps",
				HexColor: "#00FF00",
			},
			mockError: errors.New("duplicate key error"),
			wantError: true,
			errorMsg:  "error creating tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTagRepo{
				CreateFunc: func(tag models.Tag) (models.Tag, error) {
					return tt.mockTag, tt.mockError
				},
			}
			service := NewTagServise(mockRepo)

			// Act
			result, err := service.Create(tt.tag)

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
				if result.ID != tt.mockTag.ID {
					t.Errorf("Ожидался ID = %d, получили %d", tt.mockTag.ID, result.ID)
				}
			}
		})
	}
}

// TestTagService_Update тестирует метод Update
func TestTagService_Update(t *testing.T) {
	tests := []struct {
		name      string
		tag       models.Tag
		mockTag   models.Tag
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное обновление тега",
			tag: models.Tag{
				ID:       1,
				Name:     "Backend Updated",
				HexColor: "#FF5733",
			},
			mockTag: models.Tag{
				ID:       1,
				Name:     "Backend Updated",
				HexColor: "#FF5733",
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Невалидный ID (0)",
			tag: models.Tag{
				ID:       0,
				Name:     "Backend",
				HexColor: "#FF5733",
			},
			wantError: true,
			errorMsg:  "invalid tag ID",
		},
		{
			name: "Отсутствует имя",
			tag: models.Tag{
				ID:       1,
				Name:     "",
				HexColor: "#FF5733",
			},
			wantError: true,
			errorMsg:  "tag name is required",
		},
		{
			name: "Ошибка репозитория - тег не найден",
			tag: models.Tag{
				ID:       999,
				Name:     "Backend",
				HexColor: "#FF5733",
			},
			mockError: errors.New("tag not found"),
			wantError: true,
			errorMsg:  "error updating tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTagRepo{
				UpdateFunc: func(tag models.Tag) (models.Tag, error) {
					return tt.mockTag, tt.mockError
				},
			}
			service := NewTagServise(mockRepo)

			// Act
			result, err := service.Update(tt.tag)

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
				if result.Name != tt.mockTag.Name {
					t.Errorf("Ожидалось Name = %s, получили %s", tt.mockTag.Name, result.Name)
				}
			}
		})
	}
}

// TestTagService_Delete тестирует метод Delete
func TestTagService_Delete(t *testing.T) {
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
			errorMsg:  "invalid tag ID",
		},
		{
			name:      "Ошибка репозитория",
			id:        1,
			mockError: errors.New("foreign key constraint"),
			wantError: true,
			errorMsg:  "error deleting tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockTagRepo{
				DeleteFunc: func(id int64) (int64, error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewTagServise(mockRepo)

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

// TestTagService_DeleteList тестирует метод DeleteList
func TestTagService_DeleteList(t *testing.T) {
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
			mockRepo := &MockTagRepo{
				DeleteListFunc: func(ids []int64) ([]int64, error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewTagServise(mockRepo)

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

// Вспомогательная функция для проверки содержания строки
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
