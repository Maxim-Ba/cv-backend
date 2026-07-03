package repository

import (
	"fmt"
	"testing"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


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
				Description: newNullableLocalizedText("Язык программирования Go"),
				LogoUrl:     newPgText("https://go.dev/logo.png"),
			},
			wantErr: false,
		},
		{
			name: "успешное создание технологии без описания",
			technology: models.Technology{
				Title:       "Python",
				Description: i18n.NullableLocalizedText{},
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
		Description: newNullableLocalizedText("First"),
	})
	require.NoError(t, err)

	// Пытаемся создать технологию с тем же названием
	_, err = repo.Create(models.Technology{
		Title:       "Duplicate",
		Description: newNullableLocalizedText("Second"),
	})
	require.Error(t, err, "должна быть ошибка при дублировании названия")
}

func TestTechnologyRepo_Get(t *testing.T) {
	cleanupTable(t, "technology")
	repo := NewTechnologyRepo(testDB)

	// Создаем технологию для теста
	created, err := repo.Create(models.Technology{
		Title:       "TestTech",
		Description: newNullableLocalizedText("Test Description"),
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
		Description: newNullableLocalizedText("Original Description"),
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
				Description: newNullableLocalizedText("Updated Description"),
				LogoUrl:     newPgText("https://updated.com/logo.png"),
			},
			wantErr: false,
		},
		{
			name: "обновление несуществующей технологии",
			technology: models.Technology{
				ID:          99999,
				Title:       "NonExistent",
				Description: newNullableLocalizedText("Desc"),
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
		Description: newNullableLocalizedText("Will be deleted"),
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
	tech1, err := repo.Create(models.Technology{Title: "Tech1", Description: newNullableLocalizedText("Desc1")})
	require.NoError(t, err)
	tech2, err := repo.Create(models.Technology{Title: "Tech2", Description: newNullableLocalizedText("Desc2")})
	require.NoError(t, err)
	tech3, err := repo.Create(models.Technology{Title: "Tech3", Description: newNullableLocalizedText("Desc3")})
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
		{Title: "Angular", Description: newNullableLocalizedText("Frontend framework")},
		{Title: "Docker", Description: newNullableLocalizedText("Container platform")},
		{Title: "Express", Description: newNullableLocalizedText("Node.js framework")},
		{Title: "Flask", Description: newNullableLocalizedText("Python framework")},
		{Title: "Go", Description: newNullableLocalizedText("Programming language")},
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
	_, err := repo.Create(models.Technology{Title: "Golang", Description: newNullableLocalizedText("Backend")})
	require.NoError(t, err)
	_, err = repo.Create(models.Technology{Title: "React", Description: newNullableLocalizedText("Frontend")})
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
	_, err := repo.Create(models.Technology{Title: "Zebra", Description: newNullableLocalizedText("Last")})
	require.NoError(t, err)
	_, err = repo.Create(models.Technology{Title: "Alpha", Description: newNullableLocalizedText("First")})
	require.NoError(t, err)
	_, err = repo.Create(models.Technology{Title: "Middle", Description: newNullableLocalizedText("Middle")})
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

func TestTechnologyRepo_ListWithTags_PaginationWithMultipleTagsPerTech(t *testing.T) {
	cleanupAllTables(t)

	techRepo := NewTechnologyRepo(testDB)
	tagRepo := NewTagRepo(testDB)

	tag1, err := tagRepo.Create(models.Tag{Name: newLocalizedText("Backend"), HexColor: "#111111"})
	require.NoError(t, err)
	tag2, err := tagRepo.Create(models.Tag{Name: newLocalizedText("Language"), HexColor: "#222222"})
	require.NoError(t, err)

	for i := 1; i <= 8; i++ {
		tech, createErr := techRepo.Create(models.Technology{
			Title:       fmt.Sprintf("Tech%d", i),
			Description: newNullableLocalizedText(fmt.Sprintf("Description %d", i)),
		})
		require.NoError(t, createErr)
		require.NoError(t, techRepo.SetTags(tech.ID, []int64{tag1.ID, tag2.ID}))
	}

	result, err := techRepo.ListWithTags(entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
	}, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Equal(t, 8, result.Total)
	assert.Len(t, result.Content, 8, "LIMIT должен применяться к технологиям, а не к строкам JOIN")

	for i, tech := range result.Content {
		assert.Equal(t, int64(i+1), tech.ID)
		assert.Len(t, tech.Tags, 2)
	}
}

func TestTechnologyRepo_ListWithTags_PaginationPages(t *testing.T) {
	cleanupAllTables(t)

	techRepo := NewTechnologyRepo(testDB)
	tagRepo := NewTagRepo(testDB)

	tag1, err := tagRepo.Create(models.Tag{Name: newLocalizedText("Common"), HexColor: "#abcdef"})
	require.NoError(t, err)
	tag2, err := tagRepo.Create(models.Tag{Name: newLocalizedText("Extra"), HexColor: "#fedcba"})
	require.NoError(t, err)

	for i := 1; i <= 5; i++ {
		tech, createErr := techRepo.Create(models.Technology{
			Title: fmt.Sprintf("Tech%d", i),
		})
		require.NoError(t, createErr)
		require.NoError(t, techRepo.SetTags(tech.ID, []int64{tag1.ID, tag2.ID}))
	}

	page1, err := techRepo.ListWithTags(entityreqdecorator.PagebleRq{Page: 1, Size: 2}, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Equal(t, 5, page1.Total)
	assert.Len(t, page1.Content, 2)
	assert.Equal(t, int64(1), page1.Content[0].ID)
	assert.Equal(t, int64(2), page1.Content[1].ID)

	page2, err := techRepo.ListWithTags(entityreqdecorator.PagebleRq{Page: 2, Size: 2}, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Equal(t, 5, page2.Total)
	assert.Len(t, page2.Content, 2)
	assert.Equal(t, int64(3), page2.Content[0].ID)
	assert.Equal(t, int64(4), page2.Content[1].ID)

	page3, err := techRepo.ListWithTags(entityreqdecorator.PagebleRq{Page: 3, Size: 2}, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Equal(t, 5, page3.Total)
	assert.Len(t, page3.Content, 1)
	assert.Equal(t, int64(5), page3.Content[0].ID)
}

func TestTechnologyRepo_ListWithTags_SizeZeroReturnsAll(t *testing.T) {
	cleanupAllTables(t)

	techRepo := NewTechnologyRepo(testDB)
	tagRepo := NewTagRepo(testDB)

	tag, err := tagRepo.Create(models.Tag{Name: newLocalizedText("Tag"), HexColor: "#000000"})
	require.NoError(t, err)

	for i := 1; i <= 3; i++ {
		tech, createErr := techRepo.Create(models.Technology{Title: fmt.Sprintf("Tech%d", i)})
		require.NoError(t, createErr)
		require.NoError(t, techRepo.SetTags(tech.ID, []int64{tag.ID}))
	}

	result, err := techRepo.ListWithTags(entityreqdecorator.PagebleRq{Page: 1, Size: 0}, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Equal(t, 3, result.Total)
	assert.Len(t, result.Content, 3)
}

func TestTechnologyRepo_SetTags(t *testing.T) {
	cleanupTable(t, "technologies_tag")
	cleanupTable(t, "technology")
	cleanupTable(t, "tag")

	techRepo := NewTechnologyRepo(testDB)
	tagRepo := NewTagRepo(testDB)

	tech, err := techRepo.Create(models.Technology{Title: "Go"})
	require.NoError(t, err)

	tag1, err := tagRepo.Create(models.Tag{Name: newLocalizedText("Backend"), HexColor: "#111111"})
	require.NoError(t, err)
	tag2, err := tagRepo.Create(models.Tag{Name: newLocalizedText("Language"), HexColor: "#222222"})
	require.NoError(t, err)

	err = techRepo.SetTags(tech.ID, []int64{tag1.ID, tag2.ID})
	require.NoError(t, err)

	withTags, err := techRepo.GetWithTags(tech.ID, i18n.LocaleRU)
	require.NoError(t, err)
	require.Len(t, withTags.Tags, 2)
	assert.Equal(t, tag1.ID, withTags.Tags[0].ID)
	assert.Equal(t, tag2.ID, withTags.Tags[1].ID)

	err = techRepo.SetTags(tech.ID, []int64{tag1.ID})
	require.NoError(t, err)

	withTags, err = techRepo.GetWithTags(tech.ID, i18n.LocaleRU)
	require.NoError(t, err)
	require.Len(t, withTags.Tags, 1)
	assert.Equal(t, tag1.ID, withTags.Tags[0].ID)

	err = techRepo.SetTags(tech.ID, nil)
	require.NoError(t, err)

	withTags, err = techRepo.GetWithTags(tech.ID, i18n.LocaleRU)
	require.NoError(t, err)
	assert.Empty(t, withTags.Tags)
}
