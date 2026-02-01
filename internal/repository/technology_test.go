package repository

import (
	"testing"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newPgText создает pgtype.Text со значением
func newPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func TestTechnologyRepo_Create(t *testing.T) {
	cleanupTable(t, "technology")
	repo := NewTechnologyRepo(testDB)

	tests := []struct {
		name       string
		technology models.Technology
		wantErr    bool
	}{
		{
			name: "успешное создание технологии",
			technology: models.Technology{
				Title:       "Go",
				Description: newPgText("Язык программирования Go"),
				LogoUrl:     newPgText("https://go.dev/logo.png"),
			},
			wantErr: false,
		},
		{
			name: "успешное создание технологии без описания",
			technology: models.Technology{
				Title:       "Python",
				Description: pgtype.Text{Valid: false},
				LogoUrl:     pgtype.Text{Valid: false},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := repo.Create(tt.technology)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, created.ID)
			assert.Equal(t, tt.technology.Title, created.Title)
			assert.Equal(t, tt.technology.Description, created.Description)
			assert.Equal(t, tt.technology.LogoUrl, created.LogoUrl)
		})
	}
}

func TestTechnologyRepo_Create_DuplicateTitle(t *testing.T) {
	cleanupTable(t, "technology")
	repo := NewTechnologyRepo(testDB)

	// Создаем первую технологию
	_, err := repo.Create(models.Technology{
		Title:       "Duplicate",
		Description: newPgText("First"),
	})
	require.NoError(t, err)

	// Пытаемся создать технологию с тем же названием
	_, err = repo.Create(models.Technology{
		Title:       "Duplicate",
		Description: newPgText("Second"),
	})
	require.Error(t, err, "должна быть ошибка при дублировании названия")
}

func TestTechnologyRepo_Get(t *testing.T) {
	cleanupTable(t, "technology")
	repo := NewTechnologyRepo(testDB)

	// Создаем технологию для теста
	created, err := repo.Create(models.Technology{
		Title:       "TestTech",
		Description: newPgText("Test Description"),
		LogoUrl:     newPgText("https://example.com/logo.png"),
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "получение существующей технологии",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "получение несуществующей технологии",
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
			assert.Equal(t, created.Title, got.Title)
			assert.Equal(t, created.Description, got.Description)
			assert.Equal(t, created.LogoUrl, got.LogoUrl)
		})
	}
}

func TestTechnologyRepo_Update(t *testing.T) {
	cleanupTable(t, "technology")
	repo := NewTechnologyRepo(testDB)

	// Создаем технологию для теста
	created, err := repo.Create(models.Technology{
		Title:       "Original",
		Description: newPgText("Original Description"),
		LogoUrl:     newPgText("https://original.com/logo.png"),
	})
	require.NoError(t, err)

	tests := []struct {
		name       string
		technology models.Technology
		wantErr    bool
	}{
		{
			name: "успешное обновление технологии",
			technology: models.Technology{
				ID:          created.ID,
				Title:       "Updated",
				Description: newPgText("Updated Description"),
				LogoUrl:     newPgText("https://updated.com/logo.png"),
			},
			wantErr: false,
		},
		{
			name: "обновление несуществующей технологии",
			technology: models.Technology{
				ID:          99999,
				Title:       "NonExistent",
				Description: newPgText("Desc"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := repo.Update(tt.technology)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.technology.ID, updated.ID)
			assert.Equal(t, tt.technology.Title, updated.Title)
			assert.Equal(t, tt.technology.Description, updated.Description)
			assert.Equal(t, tt.technology.LogoUrl, updated.LogoUrl)
		})
	}
}

func TestTechnologyRepo_Delete(t *testing.T) {
	cleanupTable(t, "technology")
	repo := NewTechnologyRepo(testDB)

	// Создаем технологию для удаления
	created, err := repo.Create(models.Technology{
		Title:       "ToDelete",
		Description: newPgText("Will be deleted"),
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "успешное удаление технологии",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "удаление несуществующей технологии",
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

			// Проверяем, что технология действительно удалена
			_, err = repo.Get(tt.id)
			require.Error(t, err)
		})
	}
}

func TestTechnologyRepo_DeleteList(t *testing.T) {
	cleanupTable(t, "technology")
	repo := NewTechnologyRepo(testDB)

	// Создаем несколько технологий
	tech1, err := repo.Create(models.Technology{Title: "Tech1", Description: newPgText("Desc1")})
	require.NoError(t, err)
	tech2, err := repo.Create(models.Technology{Title: "Tech2", Description: newPgText("Desc2")})
	require.NoError(t, err)
	tech3, err := repo.Create(models.Technology{Title: "Tech3", Description: newPgText("Desc3")})
	require.NoError(t, err)

	tests := []struct {
		name        string
		ids         []int64
		wantDeleted int
	}{
		{
			name:        "удаление нескольких технологий",
			ids:         []int64{tech1.ID, tech2.ID},
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

	// Проверяем, что tech3 все еще существует
	got, err := repo.Get(tech3.ID)
	require.NoError(t, err)
	assert.Equal(t, tech3.Title, got.Title)
}

func TestTechnologyRepo_List(t *testing.T) {
	cleanupTable(t, "technology")
	repo := NewTechnologyRepo(testDB)

	// Создаем тестовые данные
	technologies := []models.Technology{
		{Title: "Angular", Description: newPgText("Frontend framework")},
		{Title: "Docker", Description: newPgText("Container platform")},
		{Title: "Express", Description: newPgText("Node.js framework")},
		{Title: "Flask", Description: newPgText("Python framework")},
		{Title: "Go", Description: newPgText("Programming language")},
	}

	for _, tech := range technologies {
		_, err := repo.Create(tech)
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

func TestTechnologyRepo_List_WithFilter(t *testing.T) {
	cleanupTable(t, "technology")
	repo := NewTechnologyRepo(testDB)

	// Создаем тестовые данные
	_, err := repo.Create(models.Technology{Title: "Golang", Description: newPgText("Backend")})
	require.NoError(t, err)
	_, err = repo.Create(models.Technology{Title: "React", Description: newPgText("Frontend")})
	require.NoError(t, err)

	// Фильтрация по названию
	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
		Filter: map[string]entityreqdecorator.SQLGenerator{
			"title": &entityreqdecorator.PredicateLike{
				Predicate: entityreqdecorator.Predicate{Value: "Go"},
			},
		},
	}

	result, err := repo.List(req)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Content, 1)
	assert.Equal(t, "Golang", result.Content[0].Title)
}

func TestTechnologyRepo_List_Sorting(t *testing.T) {
	cleanupTable(t, "technology")
	repo := NewTechnologyRepo(testDB)

	// Создаем тестовые данные в определенном порядке
	_, err := repo.Create(models.Technology{Title: "Zebra", Description: newPgText("Last")})
	require.NoError(t, err)
	_, err = repo.Create(models.Technology{Title: "Alpha", Description: newPgText("First")})
	require.NoError(t, err)
	_, err = repo.Create(models.Technology{Title: "Middle", Description: newPgText("Middle")})
	require.NoError(t, err)

	// Сортировка по названию ASC
	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
		Sort: []entityreqdecorator.SortBy{
			{Field: "title", Order: "ASC"},
		},
	}

	result, err := repo.List(req)
	require.NoError(t, err)
	require.Len(t, result.Content, 3)
	assert.Equal(t, "Alpha", result.Content[0].Title)
	assert.Equal(t, "Middle", result.Content[1].Title)
	assert.Equal(t, "Zebra", result.Content[2].Title)

	// Сортировка по названию DESC
	req.Sort[0].Order = "DESC"
	result, err = repo.List(req)
	require.NoError(t, err)
	require.Len(t, result.Content, 3)
	assert.Equal(t, "Zebra", result.Content[0].Title)
	assert.Equal(t, "Middle", result.Content[1].Title)
	assert.Equal(t, "Alpha", result.Content[2].Title)
}
