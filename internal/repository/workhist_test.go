package repository

import (
	"testing"
	"time"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newPgDate создает pgtype.Date со значением
func newPgDate(year int, month time.Month, day int) pgtype.Date {
	return pgtype.Date{
		Time:  time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}
}

// assertDatesEqual сравнивает pgtype.Date по значению времени, игнорируя различия в Location
func assertDatesEqual(t *testing.T, expected, actual pgtype.Date, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, expected.Valid, actual.Valid, msgAndArgs...)
	if expected.Valid && actual.Valid {
		assert.True(t, expected.Time.Equal(actual.Time), "dates should be equal: expected %v, got %v", expected.Time, actual.Time)
	}
}

func TestWorkHistoryRepo_Create(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	tests := []struct {
		name        string
		workHistory models.WorkHistory
		wantErr     bool
	}{
		{
			name: "успешное создание записи истории работы",
			workHistory: models.WorkHistory{
				Name:        "Company A",
				About:       "IT компания",
				LogoUrl:     []byte(`{"url": "https://example.com/logo.png"}`),
				PeriodStart: newPgDate(2020, time.January, 1),
				PeriodEnd:   newPgDate(2023, time.December, 31),
				WhatIDid:    []string{"Backend development", "Code review"},
				Projects:    []string{"Project A", "Project B"},
			},
			wantErr: false,
		},
		{
			name: "успешное создание записи без даты окончания",
			workHistory: models.WorkHistory{
				Name:        "Company B",
				About:       "Startup",
				PeriodStart: newPgDate(2024, time.January, 1),
				PeriodEnd:   pgtype.Date{Valid: false},
				WhatIDid:    []string{"Full stack development"},
				Projects:    []string{},
			},
			wantErr: false,
		},
		{
			name: "успешное создание записи с пустыми массивами",
			workHistory: models.WorkHistory{
				Name:        "Company C",
				About:       "Another company",
				PeriodStart: newPgDate(2019, time.June, 15),
				PeriodEnd:   newPgDate(2020, time.June, 15),
				WhatIDid:    []string{},
				Projects:    []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := repo.Create(tt.workHistory)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, created.ID)
			assert.Equal(t, tt.workHistory.Name, created.Name)
			assert.Equal(t, tt.workHistory.About, created.About)
			assertDatesEqual(t, tt.workHistory.PeriodStart, created.PeriodStart)
			assertDatesEqual(t, tt.workHistory.PeriodEnd, created.PeriodEnd)
			assert.Equal(t, tt.workHistory.WhatIDid, created.WhatIDid)
			assert.Equal(t, tt.workHistory.Projects, created.Projects)
		})
	}
}

func TestWorkHistoryRepo_Get(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Создаем запись для теста
	created, err := repo.Create(models.WorkHistory{
		Name:        "Test Company",
		About:       "Test Description",
		PeriodStart: newPgDate(2020, time.January, 1),
		PeriodEnd:   newPgDate(2023, time.December, 31),
		WhatIDid:    []string{"Task 1", "Task 2"},
		Projects:    []string{"Project 1"},
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "получение существующей записи",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "получение несуществующей записи",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, created.ID, got.ID)
			assert.Equal(t, created.Name, got.Name)
			assert.Equal(t, created.About, got.About)
			assertDatesEqual(t, created.PeriodStart, got.PeriodStart)
			assertDatesEqual(t, created.PeriodEnd, got.PeriodEnd)
			assert.Equal(t, created.WhatIDid, got.WhatIDid)
			assert.Equal(t, created.Projects, got.Projects)
		})
	}
}

func TestWorkHistoryRepo_Update(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Создаем запись для теста
	created, err := repo.Create(models.WorkHistory{
		Name:        "Original Company",
		About:       "Original About",
		PeriodStart: newPgDate(2020, time.January, 1),
		PeriodEnd:   newPgDate(2022, time.December, 31),
		WhatIDid:    []string{"Original Task"},
		Projects:    []string{"Original Project"},
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		workHistory models.WorkHistory
		wantErr     bool
	}{
		{
			name: "успешное обновление записи",
			workHistory: models.WorkHistory{
				ID:          created.ID,
				Name:        "Updated Company",
				About:       "Updated About",
				LogoUrl:     []byte(`{"url": "https://updated.com/logo.png"}`),
				PeriodStart: newPgDate(2021, time.February, 1),
				PeriodEnd:   newPgDate(2024, time.January, 15),
				WhatIDid:    []string{"Updated Task 1", "Updated Task 2"},
				Projects:    []string{"Updated Project"},
			},
			wantErr: false,
		},
		{
			name: "обновление несуществующей записи",
			workHistory: models.WorkHistory{
				ID:          99999,
				Name:        "NonExistent",
				About:       "NonExistent",
				PeriodStart: newPgDate(2020, time.January, 1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := repo.Update(tt.workHistory)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workHistory.ID, updated.ID)
			assert.Equal(t, tt.workHistory.Name, updated.Name)
			assert.Equal(t, tt.workHistory.About, updated.About)
			assertDatesEqual(t, tt.workHistory.PeriodStart, updated.PeriodStart)
			assertDatesEqual(t, tt.workHistory.PeriodEnd, updated.PeriodEnd)
			assert.Equal(t, tt.workHistory.WhatIDid, updated.WhatIDid)
			assert.Equal(t, tt.workHistory.Projects, updated.Projects)
		})
	}
}

func TestWorkHistoryRepo_Delete(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Создаем запись для удаления
	created, err := repo.Create(models.WorkHistory{
		Name:        "ToDelete Company",
		About:       "Will be deleted",
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid:    []string{},
		Projects:    []string{},
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "успешное удаление записи",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "удаление несуществующей записи",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deletedID, err := repo.Delete(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.id, deletedID)

			// Проверяем, что запись действительно удалена
			_, err = repo.Get(tt.id)
			require.Error(t, err)
		})
	}
}

func TestWorkHistoryRepo_DeleteList(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Создаем несколько записей
	wh1, err := repo.Create(models.WorkHistory{
		Name:        "Company 1",
		About:       "About 1",
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid:    []string{},
		Projects:    []string{},
	})
	require.NoError(t, err)

	wh2, err := repo.Create(models.WorkHistory{
		Name:        "Company 2",
		About:       "About 2",
		PeriodStart: newPgDate(2021, time.January, 1),
		WhatIDid:    []string{},
		Projects:    []string{},
	})
	require.NoError(t, err)

	wh3, err := repo.Create(models.WorkHistory{
		Name:        "Company 3",
		About:       "About 3",
		PeriodStart: newPgDate(2022, time.January, 1),
		WhatIDid:    []string{},
		Projects:    []string{},
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		ids         []int64
		wantDeleted int
	}{
		{
			name:        "удаление нескольких записей",
			ids:         []int64{wh1.ID, wh2.ID},
			wantDeleted: 2,
		},
		{
			name:        "удаление пустого списка",
			ids:         []int64{},
			wantDeleted: 0,
		},
		{
			name:        "удаление с несуществующими ID",
			ids:         []int64{99998, 99999},
			wantDeleted: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deletedIDs, err := repo.DeleteList(tt.ids)
			require.NoError(t, err)
			assert.Len(t, deletedIDs, tt.wantDeleted)
		})
	}

	// Проверяем, что wh3 все еще существует
	got, err := repo.Get(wh3.ID)
	require.NoError(t, err)
	assert.Equal(t, wh3.Name, got.Name)
}

func TestWorkHistoryRepo_List(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Создаем тестовые данные
	workHistories := []models.WorkHistory{
		{
			Name:        "Company A",
			About:       "About A",
			PeriodStart: newPgDate(2020, time.January, 1),
			WhatIDid:    []string{"Task A"},
			Projects:    []string{},
		},
		{
			Name:        "Company B",
			About:       "About B",
			PeriodStart: newPgDate(2021, time.February, 1),
			WhatIDid:    []string{"Task B"},
			Projects:    []string{},
		},
		{
			Name:        "Company C",
			About:       "About C",
			PeriodStart: newPgDate(2022, time.March, 1),
			WhatIDid:    []string{"Task C"},
			Projects:    []string{},
		},
		{
			Name:        "Company D",
			About:       "About D",
			PeriodStart: newPgDate(2023, time.April, 1),
			WhatIDid:    []string{"Task D"},
			Projects:    []string{},
		},
		{
			Name:        "Company E",
			About:       "About E",
			PeriodStart: newPgDate(2024, time.May, 1),
			WhatIDid:    []string{"Task E"},
			Projects:    []string{},
		},
	}

	for _, wh := range workHistories {
		_, err := repo.Create(wh)
		require.NoError(t, err)
	}

	tests := []struct {
		name        string
		req         entityreqdecorator.PagebleRq
		wantTotal   int
		wantContent int
	}{
		{
			name: "получение первой страницы",
			req: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 2,
			},
			wantTotal:   5,
			wantContent: 2,
		},
		{
			name: "получение второй страницы",
			req: entityreqdecorator.PagebleRq{
				Page: 2,
				Size: 2,
			},
			wantTotal:   5,
			wantContent: 2,
		},
		{
			name: "получение последней страницы",
			req: entityreqdecorator.PagebleRq{
				Page: 3,
				Size: 2,
			},
			wantTotal:   5,
			wantContent: 1,
		},
		{
			name: "получение всех элементов",
			req: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 10,
			},
			wantTotal:   5,
			wantContent: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.List(tt.req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantTotal, result.Total)
			assert.Len(t, result.Content, tt.wantContent)
			assert.Equal(t, tt.req.Page, result.Page)
			assert.Equal(t, tt.req.Size, result.Size)
		})
	}
}

func TestWorkHistoryRepo_List_WithFilter(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Создаем тестовые данные
	_, err := repo.Create(models.WorkHistory{
		Name:        "Yandex",
		About:       "Russian IT company",
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid:    []string{},
		Projects:    []string{},
	})
	require.NoError(t, err)

	_, err = repo.Create(models.WorkHistory{
		Name:        "Google",
		About:       "American IT company",
		PeriodStart: newPgDate(2022, time.June, 1),
		WhatIDid:    []string{},
		Projects:    []string{},
	})
	require.NoError(t, err)

	// Фильтрация по имени
	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
		Filter: map[string]entityreqdecorator.SQLGenerator{
			"name": &entityreqdecorator.PredicateLike{
				Predicate: entityreqdecorator.Predicate{Value: "Yandex"},
			},
		},
	}

	result, err := repo.List(req)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Content, 1)
	assert.Equal(t, "Yandex", result.Content[0].Name)
}

func TestWorkHistoryRepo_List_Sorting(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Создаем тестовые данные
	_, err := repo.Create(models.WorkHistory{
		Name:        "Company C",
		About:       "About C",
		PeriodStart: newPgDate(2022, time.January, 1),
		WhatIDid:    []string{},
		Projects:    []string{},
	})
	require.NoError(t, err)

	_, err = repo.Create(models.WorkHistory{
		Name:        "Company A",
		About:       "About A",
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid:    []string{},
		Projects:    []string{},
	})
	require.NoError(t, err)

	_, err = repo.Create(models.WorkHistory{
		Name:        "Company B",
		About:       "About B",
		PeriodStart: newPgDate(2021, time.January, 1),
		WhatIDid:    []string{},
		Projects:    []string{},
	})
	require.NoError(t, err)

	// Сортировка по имени ASC
	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
		Sort: []entityreqdecorator.SortBy{
			{Field: "name", Order: "ASC"},
		},
	}

	result, err := repo.List(req)
	require.NoError(t, err)
	require.Len(t, result.Content, 3)
	assert.Equal(t, "Company A", result.Content[0].Name)
	assert.Equal(t, "Company B", result.Content[1].Name)
	assert.Equal(t, "Company C", result.Content[2].Name)

	// Сортировка по дате начала DESC
	req.Sort = []entityreqdecorator.SortBy{
		{Field: "period_start", Order: "DESC"},
	}
	result, err = repo.List(req)
	require.NoError(t, err)
	require.Len(t, result.Content, 3)
	assert.Equal(t, "Company C", result.Content[0].Name) // 2022
	assert.Equal(t, "Company B", result.Content[1].Name) // 2021
	assert.Equal(t, "Company A", result.Content[2].Name) // 2020
}

func TestWorkHistoryRepo_ArrayFields(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Тестируем работу с массивами
	wh := models.WorkHistory{
		Name:        "Array Test Company",
		About:       "Testing arrays",
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid: []string{
			"Разработка микросервисов",
			"Code review",
			"Менторинг джуниоров",
			"Написание документации",
		},
		Projects: []string{
			"Проект API Gateway",
			"Проект миграции на Kubernetes",
			"Внутренний инструмент мониторинга",
		},
	}

	created, err := repo.Create(wh)
	require.NoError(t, err)
	assert.Len(t, created.WhatIDid, 4)
	assert.Len(t, created.Projects, 3)

	// Проверяем получение
	got, err := repo.Get(created.ID)
	require.NoError(t, err)
	assert.Equal(t, wh.WhatIDid, got.WhatIDid)
	assert.Equal(t, wh.Projects, got.Projects)

	// Проверяем обновление массивов
	created.WhatIDid = []string{"New task 1", "New task 2"}
	created.Projects = []string{"New project"}

	updated, err := repo.Update(created)
	require.NoError(t, err)
	assert.Len(t, updated.WhatIDid, 2)
	assert.Len(t, updated.Projects, 1)
	assert.Equal(t, "New task 1", updated.WhatIDid[0])
	assert.Equal(t, "New project", updated.Projects[0])
}
