package repository

import (
	"fmt"
	"testing"
	"time"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
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
			name: "успешное создание записи истории работы",			workHistory: models.WorkHistory{
				Name: newLocalizedText("Company A"),
				About: newLocalizedText("IT компания"),
				LogoUrl:     pgtype.Text{String: "https://example.com/logo.png", Valid: true},
				PeriodStart: newPgDate(2020, time.January, 1),
				PeriodEnd:   newPgDate(2023, time.December, 31),
				WhatIDid: newLocalizedStringList("Backend development", "Code review"),
				Projects: newLocalizedStringList("Project A", "Project B"),
			},
			wantErr: false,
		},
		{
			name: "успешное создание записи без даты окончания",			workHistory: models.WorkHistory{
				Name: newLocalizedText("Company B"),
				About: newLocalizedText("Startup"),
				PeriodStart: newPgDate(2024, time.January, 1),
				PeriodEnd:   pgtype.Date{Valid: false},
				WhatIDid: newLocalizedStringList("Full stack development"),
				Projects: newLocalizedStringList(),
			},
			wantErr: false,
		},
		{
			name: "успешное создание записи с пустыми массивами",			workHistory: models.WorkHistory{
				Name: newLocalizedText("Company C"),
				About: newLocalizedText("Another company"),
				PeriodStart: newPgDate(2019, time.June, 15),
				PeriodEnd:   newPgDate(2020, time.June, 15),
				WhatIDid: newLocalizedStringList(),
				Projects: newLocalizedStringList(),
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
		Name: newLocalizedText("Test Company"),
		About: newLocalizedText("Test Description"),
		PeriodStart: newPgDate(2020, time.January, 1),
		PeriodEnd:   newPgDate(2023, time.December, 31),
		WhatIDid: newLocalizedStringList("Task 1", "Task 2"),
		Projects: newLocalizedStringList("Project 1"),
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name: "получение существующей записи",			id:      created.ID,
			wantErr: false,
		},
		{
			name: "получение несуществующей записи",			id:      99999,
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
		Name: newLocalizedText("Original Company"),
		About: newLocalizedText("Original About"),
		PeriodStart: newPgDate(2020, time.January, 1),
		PeriodEnd:   newPgDate(2022, time.December, 31),
		WhatIDid: newLocalizedStringList("Original Task"),
		Projects: newLocalizedStringList("Original Project"),
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		workHistory models.WorkHistory
		wantErr     bool
	}{
		{
			name: "успешное обновление записи",			workHistory: models.WorkHistory{
				ID:          created.ID,
				Name: newLocalizedText("Updated Company"),
				About: newLocalizedText("Updated About"),
				LogoUrl:     pgtype.Text{String: "https://updated.com/logo.png", Valid: true},
				PeriodStart: newPgDate(2021, time.February, 1),
				PeriodEnd:   newPgDate(2024, time.January, 15),
				WhatIDid: newLocalizedStringList("Updated Task 1", "Updated Task 2"),
				Projects: newLocalizedStringList("Updated Project"),
			},
			wantErr: false,
		},
		{
			name: "обновление несуществующей записи",			workHistory: models.WorkHistory{
				ID:          99999,
				Name: newLocalizedText("NonExistent"),
				About: newLocalizedText("NonExistent"),
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
		Name: newLocalizedText("ToDelete Company"),
		About: newLocalizedText("Will be deleted"),
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name: "успешное удаление записи",			id:      created.ID,
			wantErr: false,
		},
		{
			name: "удаление несуществующей записи",			id:      99999,
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
		Name: newLocalizedText("Company 1"),
		About: newLocalizedText("About 1"),
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	wh2, err := repo.Create(models.WorkHistory{
		Name: newLocalizedText("Company 2"),
		About: newLocalizedText("About 2"),
		PeriodStart: newPgDate(2021, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	wh3, err := repo.Create(models.WorkHistory{
		Name: newLocalizedText("Company 3"),
		About: newLocalizedText("About 3"),
		PeriodStart: newPgDate(2022, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		ids         []int64
		wantDeleted int
	}{
		{
			name: "удаление нескольких записей",			ids:         []int64{wh1.ID, wh2.ID},
			wantDeleted: 2,
		},
		{
			name: "удаление пустого списка",			ids:         []int64{},
			wantDeleted: 0,
		},
		{
			name: "удаление с несуществующими ID",			ids:         []int64{99998, 99999},
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
			Name: newLocalizedText("Company A"),
			About: newLocalizedText("About A"),
			PeriodStart: newPgDate(2020, time.January, 1),
			WhatIDid: newLocalizedStringList("Task A"),
			Projects: newLocalizedStringList(),
		},
		{
			Name: newLocalizedText("Company B"),
			About: newLocalizedText("About B"),
			PeriodStart: newPgDate(2021, time.February, 1),
			WhatIDid: newLocalizedStringList("Task B"),
			Projects: newLocalizedStringList(),
		},
		{
			Name: newLocalizedText("Company C"),
			About: newLocalizedText("About C"),
			PeriodStart: newPgDate(2022, time.March, 1),
			WhatIDid: newLocalizedStringList("Task C"),
			Projects: newLocalizedStringList(),
		},
		{
			Name: newLocalizedText("Company D"),
			About: newLocalizedText("About D"),
			PeriodStart: newPgDate(2023, time.April, 1),
			WhatIDid: newLocalizedStringList("Task D"),
			Projects: newLocalizedStringList(),
		},
		{
			Name: newLocalizedText("Company E"),
			About: newLocalizedText("About E"),
			PeriodStart: newPgDate(2024, time.May, 1),
			WhatIDid: newLocalizedStringList("Task E"),
			Projects: newLocalizedStringList(),
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
	created, err := repo.Create(models.WorkHistory{
		Name: newLocalizedText("Yandex"),
		About: newLocalizedText("Russian IT company"),
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	_, err = repo.Create(models.WorkHistory{
		Name: newLocalizedText("Google"),
		About: newLocalizedText("American IT company"),
		PeriodStart: newPgDate(2022, time.June, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
		Filter: map[string]entityreqdecorator.SQLGenerator{
			"id": &entityreqdecorator.PredicateEQ{
				Predicate: entityreqdecorator.Predicate{Value: fmt.Sprintf("%d", created.ID)},
			},
		},
	}

	result, err := repo.List(req)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Content, 1)
	assert.Equal(t, "Yandex", result.Content[0].Name.Get(i18n.LocaleRU))
}

func TestWorkHistoryRepo_List_Sorting(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Создаем тестовые данные
	_, err := repo.Create(models.WorkHistory{
		Name: newLocalizedText("Company C"),
		About: newLocalizedText("About C"),
		PeriodStart: newPgDate(2022, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	_, err = repo.Create(models.WorkHistory{
		Name: newLocalizedText("Company A"),
		About: newLocalizedText("About A"),
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	_, err = repo.Create(models.WorkHistory{
		Name: newLocalizedText("Company B"),
		About: newLocalizedText("About B"),
		PeriodStart: newPgDate(2021, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
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
	assert.Equal(t, "Company A", result.Content[0].Name.Get(i18n.LocaleRU))
	assert.Equal(t, "Company B", result.Content[1].Name.Get(i18n.LocaleRU))
	assert.Equal(t, "Company C", result.Content[2].Name.Get(i18n.LocaleRU))

	// Сортировка по дате начала DESC
	req.Sort = []entityreqdecorator.SortBy{
		{Field: "period_start", Order: "DESC"},
	}
	result, err = repo.List(req)
	require.NoError(t, err)
	require.Len(t, result.Content, 3)
	assert.Equal(t, "Company C", result.Content[0].Name.Get(i18n.LocaleRU)) // 2022
	assert.Equal(t, "Company B", result.Content[1].Name.Get(i18n.LocaleRU)) // 2021
	assert.Equal(t, "Company A", result.Content[2].Name.Get(i18n.LocaleRU)) // 2020
}

func TestWorkHistoryRepo_ArrayFields(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Тестируем работу с массивами
	wh := models.WorkHistory{
		Name: newLocalizedText("Array Test Company"),
		About: newLocalizedText("Testing arrays"),
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid: newLocalizedStringList(
			"Разработка микросервисов",
			"Code review",
			"Менторинг джуниоров",
			"Написание документации",
		),
		Projects: newLocalizedStringList(
			"Проект API Gateway",
			"Проект миграции на Kubernetes",
			"Внутренний инструмент мониторинга",
		),
	}

	created, err := repo.Create(wh)
	require.NoError(t, err)
	assert.Len(t, created.WhatIDid.Get(i18n.LocaleRU), 4)
	assert.Len(t, created.Projects.Get(i18n.LocaleRU), 3)

	// Проверяем получение
	got, err := repo.Get(created.ID)
	require.NoError(t, err)
	assert.Equal(t, wh.WhatIDid, got.WhatIDid)
	assert.Equal(t, wh.Projects, got.Projects)

	// Проверяем обновление массивов
	created.WhatIDid = newLocalizedStringList("New task 1", "New task 2")
	created.Projects = newLocalizedStringList("New project")

	updated, err := repo.Update(created)
	require.NoError(t, err)
	assert.Len(t, updated.WhatIDid.Get(i18n.LocaleRU), 2)
	assert.Len(t, updated.Projects.Get(i18n.LocaleRU), 1)
	assert.Equal(t, "New task 1", updated.WhatIDid.Get(i18n.LocaleRU)[0])
	assert.Equal(t, "New project", updated.Projects.Get(i18n.LocaleRU)[0])
}

func TestWorkHistoryRepo_ListWithTechnologies_EmptyDatabase(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
	}

	result, err := repo.ListWithTechnologies(req, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Equal(t, 0, result.Total)
	assert.Empty(t, result.Content)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Size)
}

func TestWorkHistoryRepo_ListWithTechnologies_WithoutTechnologies(t *testing.T) {
	cleanupTable(t, "work_history")
	repo := NewWorkHistoryRepo(testDB)

	// Создаем записи истории работы без технологий
	wh1, err := repo.Create(models.WorkHistory{
		Name: newLocalizedText("Company A"),
		About: newLocalizedText("About A"),
		LogoUrl:     pgtype.Text{String: "https://example.com/a.png", Valid: true},
		PeriodStart: newPgDate(2020, time.January, 1),
		PeriodEnd:   newPgDate(2021, time.December, 31),
		WhatIDid: newLocalizedStringList("Task A1", "Task A2"),
		Projects: newLocalizedStringList("Project A"),
	})
	require.NoError(t, err)

	wh2, err := repo.Create(models.WorkHistory{
		Name: newLocalizedText("Company B"),
		About: newLocalizedText("About B"),
		PeriodStart: newPgDate(2022, time.March, 15),
		PeriodEnd:   pgtype.Date{Valid: false},
		WhatIDid: newLocalizedStringList("Task B1"),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
	}

	result, err := repo.ListWithTechnologies(req, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Content, 2)

	// Проверяем первую запись
	assert.Equal(t, wh1.ID, result.Content[0].ID)
	assert.Equal(t, "Company A", result.Content[0].Name)
	assert.Equal(t, "About A", result.Content[0].About)
	assert.Equal(t, "https://example.com/a.png", result.Content[0].LogoUrl)
	assert.NotNil(t, result.Content[0].PeriodStart)
	assert.Equal(t, "2020-01-01", *result.Content[0].PeriodStart)
	assert.NotNil(t, result.Content[0].PeriodEnd)
	assert.Equal(t, "2021-12-31", *result.Content[0].PeriodEnd)
	assert.Equal(t, []string{"Task A1", "Task A2"}, result.Content[0].WhatIDid)
	assert.Equal(t, []string{"Project A"}, result.Content[0].Projects)
	assert.Empty(t, result.Content[0].Technologies)

	// Проверяем вторую запись
	assert.Equal(t, wh2.ID, result.Content[1].ID)
	assert.Equal(t, "Company B", result.Content[1].Name)
	assert.NotNil(t, result.Content[1].PeriodStart)
	assert.Equal(t, "2022-03-15", *result.Content[1].PeriodStart)
	assert.Nil(t, result.Content[1].PeriodEnd)
	assert.Empty(t, result.Content[1].Technologies)
}

func TestWorkHistoryRepo_ListWithTechnologies_WithTechnologiesAndTags(t *testing.T) {
	cleanupTable(t, "work_history_technology")
	cleanupTable(t, "technologies_tag")
	cleanupTable(t, "work_history")
	cleanupTable(t, "technology")
	cleanupTable(t, "tag")

	whRepo := NewWorkHistoryRepo(testDB)
	techRepo := NewTechnologyRepo(testDB)
	tagRepo := NewTagRepo(testDB)

	// Создаем теги
	tag1, err := tagRepo.Create(models.Tag{
		Name: newLocalizedText("Backend"),
		HexColor: "#FF5733",
	})
	require.NoError(t, err)

	tag2, err := tagRepo.Create(models.Tag{
		Name: newLocalizedText("Database"),
		HexColor: "#33FF57",
	})
	require.NoError(t, err)

	tag3, err := tagRepo.Create(models.Tag{
		Name: newLocalizedText("Frontend"),
		HexColor: "#3357FF",
	})
	require.NoError(t, err)

	// Создаем технологии
	tech1, err := techRepo.Create(models.Technology{
		Title:       "Go",
		Description: newNullableLocalizedText("Programming language"),
		LogoUrl:     pgtype.Text{String: "https://golang.org/logo.png", Valid: true},
	})
	require.NoError(t, err)

	tech2, err := techRepo.Create(models.Technology{
		Title:       "PostgreSQL",
		Description: newNullableLocalizedText("Relational database"),
		LogoUrl:     pgtype.Text{Valid: false},
	})
	require.NoError(t, err)

	tech3, err := techRepo.Create(models.Technology{
		Title: "React",
	})
	require.NoError(t, err)

	// Связываем технологии с тегами
	_, err = testDB.Exec("INSERT INTO technologies_tag (technology_id, tag_id) VALUES ($1, $2)", tech1.ID, tag1.ID)
	require.NoError(t, err)

	_, err = testDB.Exec("INSERT INTO technologies_tag (technology_id, tag_id) VALUES ($1, $2)", tech2.ID, tag1.ID)
	require.NoError(t, err)

	_, err = testDB.Exec("INSERT INTO technologies_tag (technology_id, tag_id) VALUES ($1, $2)", tech2.ID, tag2.ID)
	require.NoError(t, err)

	_, err = testDB.Exec("INSERT INTO technologies_tag (technology_id, tag_id) VALUES ($1, $2)", tech3.ID, tag3.ID)
	require.NoError(t, err)

	// Создаем историю работы
	wh1, err := whRepo.Create(models.WorkHistory{
		Name: newLocalizedText("Tech Company"),
		About: newLocalizedText("Full stack development"),
		PeriodStart: newPgDate(2020, time.January, 1),
		PeriodEnd:   newPgDate(2023, time.December, 31),
		WhatIDid: newLocalizedStringList("Backend", "Database design"),
		Projects: newLocalizedStringList("Project X"),
	})
	require.NoError(t, err)

	wh2, err := whRepo.Create(models.WorkHistory{
		Name: newLocalizedText("Startup"),
		About: newLocalizedText("Frontend development"),
		PeriodStart: newPgDate(2024, time.January, 1),
		WhatIDid: newLocalizedStringList("UI development"),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	// Связываем историю работы с технологиями
	_, err = testDB.Exec("INSERT INTO work_history_technology (work_history_id, technology_id) VALUES ($1, $2)", wh1.ID, tech1.ID)
	require.NoError(t, err)

	_, err = testDB.Exec("INSERT INTO work_history_technology (work_history_id, technology_id) VALUES ($1, $2)", wh1.ID, tech2.ID)
	require.NoError(t, err)

	_, err = testDB.Exec("INSERT INTO work_history_technology (work_history_id, technology_id) VALUES ($1, $2)", wh2.ID, tech3.ID)
	require.NoError(t, err)

	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
	}

	result, err := whRepo.ListWithTechnologies(req, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Content, 2)

	// Проверяем первую запись (Tech Company с Go и PostgreSQL)
	assert.Equal(t, wh1.ID, result.Content[0].ID)
	assert.Equal(t, "Tech Company", result.Content[0].Name)
	assert.Len(t, result.Content[0].Technologies, 2)

	// Проверяем технологию Go
	goTech := result.Content[0].Technologies[0]
	assert.Equal(t, tech1.ID, goTech.ID)
	assert.Equal(t, "Go", goTech.Title)
	assert.NotNil(t, goTech.Description)
	assert.Equal(t, "Programming language", *goTech.Description)
	assert.NotNil(t, goTech.LogoUrl)
	assert.Equal(t, "https://golang.org/logo.png", *goTech.LogoUrl)
	assert.Len(t, goTech.Tags, 1)
	assert.Equal(t, tag1.ID, goTech.Tags[0].ID)
	assert.Equal(t, "Backend", goTech.Tags[0].Name)
	assert.Equal(t, "#FF5733", goTech.Tags[0].HexColor)

	// Проверяем технологию PostgreSQL
	pgTech := result.Content[0].Technologies[1]
	assert.Equal(t, tech2.ID, pgTech.ID)
	assert.Equal(t, "PostgreSQL", pgTech.Title)
	assert.NotNil(t, pgTech.Description)
	assert.Equal(t, "Relational database", *pgTech.Description)
	assert.Nil(t, pgTech.LogoUrl)
	assert.Len(t, pgTech.Tags, 2)

	// Проверяем вторую запись (Startup с React)
	assert.Equal(t, wh2.ID, result.Content[1].ID)
	assert.Equal(t, "Startup", result.Content[1].Name)
	assert.Len(t, result.Content[1].Technologies, 1)

	reactTech := result.Content[1].Technologies[0]
	assert.Equal(t, tech3.ID, reactTech.ID)
	assert.Equal(t, "React", reactTech.Title)
	assert.Nil(t, reactTech.Description)
	assert.Nil(t, reactTech.LogoUrl)
	assert.Len(t, reactTech.Tags, 1)
	assert.Equal(t, tag3.ID, reactTech.Tags[0].ID)
	assert.Equal(t, "Frontend", reactTech.Tags[0].Name)
}

func TestWorkHistoryRepo_ListWithTechnologies_Pagination(t *testing.T) {
	cleanupTable(t, "work_history_technology")
	cleanupTable(t, "work_history")

	repo := NewWorkHistoryRepo(testDB)

	// Создаем 5 записей истории работы
	for i := 1; i <= 5; i++ {
		_, err := repo.Create(models.WorkHistory{
			Name: newLocalizedText(fmt.Sprintf("Company %d", i)),
			About: newLocalizedText(fmt.Sprintf("About %d", i)),
			PeriodStart: newPgDate(2020+i-1, time.January, 1),
			WhatIDid: newLocalizedStringList(),
			Projects: newLocalizedStringList(),
		})
		require.NoError(t, err)
	}

	tests := []struct {
		name        string
		req         entityreqdecorator.PagebleRq
		wantTotal   int
		wantContent int
		checkFirst  string
	}{
		{
			name: "первая страница с размером 2",
			req: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 2,
			},
			wantTotal:   5,
			wantContent: 2,
			checkFirst:  "Company 1",
		},
		{
			name: "вторая страница с размером 2",
			req: entityreqdecorator.PagebleRq{
				Page: 2,
				Size: 2,
			},
			wantTotal:   5,
			wantContent: 2,
			checkFirst:  "Company 3",
		},
		{
			name: "последняя страница с размером 2",
			req: entityreqdecorator.PagebleRq{
				Page: 3,
				Size: 2,
			},
			wantTotal:   5,
			wantContent: 1,
			checkFirst:  "Company 5",
		},
		{
			name: "размер больше общего количества",
			req: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 10,
			},
			wantTotal:   5,
			wantContent: 5,
			checkFirst:  "Company 1",
		},
		{
			name: "нулевой размер (получить все)",
			req: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 0,
			},
			wantTotal:   5,
			wantContent: 5,
			checkFirst:  "Company 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.ListWithTechnologies(tt.req, i18n.LocaleRU)
			require.NoError(t, err)
			assert.Equal(t, tt.wantTotal, result.Total)
			assert.Len(t, result.Content, tt.wantContent)
			assert.Equal(t, tt.req.Page, result.Page)
			assert.Equal(t, tt.req.Size, result.Size)

			if tt.wantContent > 0 {
				assert.Equal(t, tt.checkFirst, result.Content[0].Name)
			}
		})
	}
}

func TestWorkHistoryRepo_ListWithTechnologies_OrderPreservation(t *testing.T) {
	cleanupTable(t, "work_history_technology")
	cleanupTable(t, "technologies_tag")
	cleanupTable(t, "work_history")
	cleanupTable(t, "technology")
	cleanupTable(t, "tag")

	whRepo := NewWorkHistoryRepo(testDB)
	techRepo := NewTechnologyRepo(testDB)

	// Создаем технологии
	tech1, err := techRepo.Create(models.Technology{Title: "Tech A"})
	require.NoError(t, err)

	tech2, err := techRepo.Create(models.Technology{Title: "Tech B"})
	require.NoError(t, err)

	tech3, err := techRepo.Create(models.Technology{Title: "Tech C"})
	require.NoError(t, err)

	// Создаем историю работы
	wh1, err := whRepo.Create(models.WorkHistory{
		Name: newLocalizedText("Company 1"),
		About: newLocalizedText("About 1"),
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	wh2, err := whRepo.Create(models.WorkHistory{
		Name: newLocalizedText("Company 2"),
		About: newLocalizedText("About 2"),
		PeriodStart: newPgDate(2021, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	wh3, err := whRepo.Create(models.WorkHistory{
		Name: newLocalizedText("Company 3"),
		About: newLocalizedText("About 3"),
		PeriodStart: newPgDate(2022, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	// Связываем в определенном порядке
	_, err = testDB.Exec("INSERT INTO work_history_technology (work_history_id, technology_id) VALUES ($1, $2)", wh1.ID, tech1.ID)
	require.NoError(t, err)

	_, err = testDB.Exec("INSERT INTO work_history_technology (work_history_id, technology_id) VALUES ($1, $2)", wh2.ID, tech2.ID)
	require.NoError(t, err)

	_, err = testDB.Exec("INSERT INTO work_history_technology (work_history_id, technology_id) VALUES ($1, $2)", wh3.ID, tech3.ID)
	require.NoError(t, err)

	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
	}

	result, err := whRepo.ListWithTechnologies(req, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Len(t, result.Content, 3)

	// Проверяем, что порядок сохраняется (по ID work_history)
	assert.Equal(t, wh1.ID, result.Content[0].ID)
	assert.Equal(t, wh2.ID, result.Content[1].ID)
	assert.Equal(t, wh3.ID, result.Content[2].ID)
}

func TestWorkHistoryRepo_ListWithTechnologies_TechnologyWithoutTags(t *testing.T) {
	cleanupTable(t, "work_history_technology")
	cleanupTable(t, "technologies_tag")
	cleanupTable(t, "work_history")
	cleanupTable(t, "technology")

	whRepo := NewWorkHistoryRepo(testDB)
	techRepo := NewTechnologyRepo(testDB)

	// Создаем технологию без тегов
	tech, err := techRepo.Create(models.Technology{
		Title:       "Standalone Tech",
		Description: newNullableLocalizedText("No tags"),
	})
	require.NoError(t, err)

	// Создаем историю работы
	wh, err := whRepo.Create(models.WorkHistory{
		Name: newLocalizedText("Company"),
		About: newLocalizedText("About"),
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	// Связываем
	_, err = testDB.Exec("INSERT INTO work_history_technology (work_history_id, technology_id) VALUES ($1, $2)", wh.ID, tech.ID)
	require.NoError(t, err)

	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
	}

	result, err := whRepo.ListWithTechnologies(req, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Len(t, result.Content, 1)
	assert.Len(t, result.Content[0].Technologies, 1)

	techResult := result.Content[0].Technologies[0]
	assert.Equal(t, tech.ID, techResult.ID)
	assert.Equal(t, "Standalone Tech", techResult.Title)
	assert.NotNil(t, techResult.Description)
	assert.Equal(t, "No tags", *techResult.Description)
	assert.Empty(t, techResult.Tags)
}

func TestWorkHistoryRepo_SetTechnologies(t *testing.T) {
	cleanupTable(t, "work_history_technology")
	cleanupTable(t, "work_history")
	cleanupTable(t, "technology")

	whRepo := NewWorkHistoryRepo(testDB)
	techRepo := NewTechnologyRepo(testDB)

	wh, err := whRepo.Create(models.WorkHistory{
		Name: newLocalizedText("Company"),
		About: newLocalizedText("About company"),
		PeriodStart: newPgDate(2020, time.January, 1),
		WhatIDid: newLocalizedStringList(),
		Projects: newLocalizedStringList(),
	})
	require.NoError(t, err)

	tech1, err := techRepo.Create(models.Technology{Title: "Go"})
	require.NoError(t, err)
	tech2, err := techRepo.Create(models.Technology{Title: "PostgreSQL"})
	require.NoError(t, err)

	err = whRepo.SetTechnologies(wh.ID, []int64{tech1.ID, tech2.ID})
	require.NoError(t, err)

	withTech, err := whRepo.GetWithTechnologies(wh.ID, i18n.LocaleRU)
	require.NoError(t, err)
	require.Len(t, withTech.Technologies, 2)
	assert.Equal(t, tech1.ID, withTech.Technologies[0].ID)
	assert.Equal(t, tech2.ID, withTech.Technologies[1].ID)

	err = whRepo.SetTechnologies(wh.ID, []int64{tech1.ID})
	require.NoError(t, err)

	withTech, err = whRepo.GetWithTechnologies(wh.ID, i18n.LocaleRU)
	require.NoError(t, err)
	require.Len(t, withTech.Technologies, 1)
	assert.Equal(t, tech1.ID, withTech.Technologies[0].ID)

	err = whRepo.SetTechnologies(wh.ID, nil)
	require.NoError(t, err)

	withTech, err = whRepo.GetWithTechnologies(wh.ID, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Empty(t, withTech.Technologies)
}
