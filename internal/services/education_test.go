package services

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

// MockEducationRepo мок-репозиторий для тестирования EducationService
type MockEducationRepo struct {
	GetFunc        func(id int64) (models.Education, error)
	ListFunc       func(entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Education], error)
	CreateFunc     func(models.Education) (models.Education, error)
	UpdateFunc     func(models.Education) (models.Education, error)
	DeleteFunc     func(id int64) (int64, error)
	DeleteListFunc func([]int64) ([]int64, error)
}

func (m *MockEducationRepo) Get(id int64) (models.Education, error) {
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return models.Education{}, nil
}

func (m *MockEducationRepo) List(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Education], error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return entityreqdecorator.PagebleRs[models.Education]{}, nil
}

func (m *MockEducationRepo) Create(edu models.Education) (models.Education, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(edu)
	}
	return models.Education{}, nil
}

func (m *MockEducationRepo) Update(edu models.Education) (models.Education, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(edu)
	}
	return models.Education{}, nil
}

func (m *MockEducationRepo) Delete(id int64) (int64, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return 0, nil
}

func (m *MockEducationRepo) DeleteList(ids []int64) ([]int64, error) {
	if m.DeleteListFunc != nil {
		return m.DeleteListFunc(ids)
	}
	return nil, nil
}

// TestEducationService_Get тестирует метод Get
func TestEducationService_Get(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		mockEdu   models.Education
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное получение образования",
			id:   1,
			mockEdu: models.Education{
				ID:           1,
				Name:         pgtype.Text{String: "МГУ", Valid: true},
				Year:         2020,
				Course:       "Computer Science",
				Organization: "Moscow State University",
			},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Невалидный ID (0)",
			id:        0,
			mockError: nil,
			wantError: true,
			errorMsg:  "invalid education ID",
		},
		{
			name:      "Ошибка репозитория",
			id:        1,
			mockError: errors.New("database error"),
			wantError: true,
			errorMsg:  "error getting education",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEducationRepo{
				GetFunc: func(id int64) (models.Education, error) {
					return tt.mockEdu, tt.mockError
				},
			}
			service := NewEducationService(mockRepo)

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
				if result.ID != tt.mockEdu.ID {
					t.Errorf("Ожидался ID = %d, получили %d", tt.mockEdu.ID, result.ID)
				}
				if result.Course != tt.mockEdu.Course {
					t.Errorf("Ожидался Course = %s, получили %s", tt.mockEdu.Course, result.Course)
				}
			}
		})
	}
}

// TestEducationService_List тестирует метод List
func TestEducationService_List(t *testing.T) {
	tests := []struct {
		name       string
		request    entityreqdecorator.PagebleRq
		mockResult entityreqdecorator.PagebleRs[models.Education]
		mockError  error
		wantError  bool
	}{
		{
			name: "Успешное получение списка образований",
			request: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 10,
			},
			mockResult: entityreqdecorator.PagebleRs[models.Education]{
				Total: 2,
				Content: []models.Education{
					{
						ID:           1,
						Name:         pgtype.Text{String: "МГУ", Valid: true},
						Year:         2020,
						Course:       "Computer Science",
						Organization: "Moscow State University",
					},
					{
						ID:           2,
						Name:         pgtype.Text{String: "МФТИ", Valid: true},
						Year:         2019,
						Course:       "Physics",
						Organization: "Moscow Institute of Physics and Technology",
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
			mockRepo := &MockEducationRepo{
				ListFunc: func(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Education], error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewEducationService(mockRepo)

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

// TestEducationService_Create тестирует метод Create
func TestEducationService_Create(t *testing.T) {
	tests := []struct {
		name      string
		edu       models.Education
		mockEdu   models.Education
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное создание образования",
			edu: models.Education{
				Name:         pgtype.Text{String: "СПбГУ", Valid: true},
				Year:         2021,
				Course:       "Mathematics",
				Organization: "Saint Petersburg State University",
			},
			mockEdu: models.Education{
				ID:           3,
				Name:         pgtype.Text{String: "СПбГУ", Valid: true},
				Year:         2021,
				Course:       "Mathematics",
				Organization: "Saint Petersburg State University",
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Отсутствует курс",
			edu: models.Education{
				Name:         pgtype.Text{String: "СПбГУ", Valid: true},
				Year:         2021,
				Course:       "",
				Organization: "Saint Petersburg State University",
			},
			wantError: true,
			errorMsg:  "course and organization are required fields",
		},
		{
			name: "Отсутствует организация",
			edu: models.Education{
				Name:         pgtype.Text{String: "СПбГУ", Valid: true},
				Year:         2021,
				Course:       "Mathematics",
				Organization: "",
			},
			wantError: true,
			errorMsg:  "course and organization are required fields",
		},
		{
			name: "Ошибка репозитория",
			edu: models.Education{
				Name:         pgtype.Text{String: "СПбГУ", Valid: true},
				Year:         2021,
				Course:       "Mathematics",
				Organization: "Saint Petersburg State University",
			},
			mockError: errors.New("database constraint violation"),
			wantError: true,
			errorMsg:  "error creating education",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEducationRepo{
				CreateFunc: func(edu models.Education) (models.Education, error) {
					return tt.mockEdu, tt.mockError
				},
			}
			service := NewEducationService(mockRepo)

			// Act
			result, err := service.Create(tt.edu)

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
				if result.ID != tt.mockEdu.ID {
					t.Errorf("Ожидался ID = %d, получили %d", tt.mockEdu.ID, result.ID)
				}
			}
		})
	}
}

// TestEducationService_Update тестирует метод Update
func TestEducationService_Update(t *testing.T) {
	tests := []struct {
		name      string
		edu       models.Education
		mockEdu   models.Education
		mockError error
		wantError bool
		errorMsg  string
	}{
		{
			name: "Успешное обновление образования",
			edu: models.Education{
				ID:           1,
				Name:         pgtype.Text{String: "МГУ", Valid: true},
				Year:         2021,
				Course:       "Computer Science Updated",
				Organization: "Moscow State University",
			},
			mockEdu: models.Education{
				ID:           1,
				Name:         pgtype.Text{String: "МГУ", Valid: true},
				Year:         2021,
				Course:       "Computer Science Updated",
				Organization: "Moscow State University",
			},
			mockError: nil,
			wantError: false,
		},
		{
			name: "Невалидный ID (0)",
			edu: models.Education{
				ID:           0,
				Course:       "Computer Science",
				Organization: "Some University",
			},
			wantError: true,
			errorMsg:  "invalid education ID",
		},
		{
			name: "Отсутствует курс",
			edu: models.Education{
				ID:           1,
				Course:       "",
				Organization: "Some University",
			},
			wantError: true,
			errorMsg:  "course and organization are required fields",
		},
		{
			name: "Ошибка репозитория - образование не найдено",
			edu: models.Education{
				ID:           999,
				Course:       "Computer Science",
				Organization: "Some University",
			},
			mockError: errors.New("education not found"),
			wantError: true,
			errorMsg:  "error updating education",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEducationRepo{
				UpdateFunc: func(edu models.Education) (models.Education, error) {
					return tt.mockEdu, tt.mockError
				},
			}
			service := NewEducationService(mockRepo)

			// Act
			result, err := service.Update(tt.edu)

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
				if result.Course != tt.mockEdu.Course {
					t.Errorf("Ожидался Course = %s, получили %s", tt.mockEdu.Course, result.Course)
				}
			}
		})
	}
}

// TestEducationService_Delete тестирует метод Delete
func TestEducationService_Delete(t *testing.T) {
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
			errorMsg:  "invalid education ID",
		},
		{
			name:      "Ошибка репозитория",
			id:        1,
			mockError: errors.New("foreign key constraint"),
			wantError: true,
			errorMsg:  "error deleting education",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockEducationRepo{
				DeleteFunc: func(id int64) (int64, error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewEducationService(mockRepo)

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

// TestEducationService_DeleteList тестирует метод DeleteList
func TestEducationService_DeleteList(t *testing.T) {
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
			mockRepo := &MockEducationRepo{
				DeleteListFunc: func(ids []int64) ([]int64, error) {
					return tt.mockResult, tt.mockError
				},
			}
			service := NewEducationService(mockRepo)

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
