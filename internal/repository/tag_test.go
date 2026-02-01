package repository

import (
	"testing"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagRepo_Create(t *testing.T) {
	cleanupTable(t, "tag")
	repo := NewTagRepo(testDB)

	tests := []struct {
		name    string
		tag     models.Tag
		wantErr bool
	}{
		{
			name: "успешное создание тега",
			tag: models.Tag{
				Name:     "Backend",
				HexColor: "#FF5733",
			},
			wantErr: false,
		},
		{
			name: "успешное создание второго тега",
			tag: models.Tag{
				Name:     "Frontend",
				HexColor: "#33FF57",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := repo.Create(tt.tag)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, created.ID)
			assert.Equal(t, tt.tag.Name, created.Name)
			assert.Equal(t, tt.tag.HexColor, created.HexColor)
		})
	}
}

func TestTagRepo_Create_DuplicateName(t *testing.T) {
	cleanupTable(t, "tag")
	repo := NewTagRepo(testDB)

	// Создаем первый тег
	_, err := repo.Create(models.Tag{
		Name:     "Duplicate",
		HexColor: "#111111",
	})
	require.NoError(t, err)

	// Пытаемся создать тег с тем же именем
	_, err = repo.Create(models.Tag{
		Name:     "Duplicate",
		HexColor: "#222222",
	})
	require.Error(t, err, "должна быть ошибка при дублировании имени")
}

func TestTagRepo_Get(t *testing.T) {
	cleanupTable(t, "tag")
	repo := NewTagRepo(testDB)

	// Создаем тег для теста
	created, err := repo.Create(models.Tag{
		Name:     "TestGet",
		HexColor: "#ABCDEF",
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "получение существующего тега",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "получение несуществующего тега",
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
			assert.Equal(t, created.HexColor, got.HexColor)
		})
	}
}

func TestTagRepo_Update(t *testing.T) {
	cleanupTable(t, "tag")
	repo := NewTagRepo(testDB)

	// Создаем тег для теста
	created, err := repo.Create(models.Tag{
		Name:     "Original",
		HexColor: "#000000",
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		tag     models.Tag
		wantErr bool
	}{
		{
			name: "успешное обновление тега",
			tag: models.Tag{
				ID:       created.ID,
				Name:     "Updated",
				HexColor: "#FFFFFF",
			},
			wantErr: false,
		},
		{
			name: "обновление несуществующего тега",
			tag: models.Tag{
				ID:       99999,
				Name:     "NonExistent",
				HexColor: "#123456",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := repo.Update(tt.tag)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.tag.ID, updated.ID)
			assert.Equal(t, tt.tag.Name, updated.Name)
			assert.Equal(t, tt.tag.HexColor, updated.HexColor)
		})
	}
}

func TestTagRepo_Delete(t *testing.T) {
	cleanupTable(t, "tag")
	repo := NewTagRepo(testDB)

	// Создаем тег для удаления
	created, err := repo.Create(models.Tag{
		Name:     "ToDelete",
		HexColor: "#AABBCC",
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "успешное удаление тега",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "удаление несуществующего тега",
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

			// Проверяем, что тег действительно удален
			_, err = repo.Get(tt.id)
			require.Error(t, err)
		})
	}
}

func TestTagRepo_DeleteList(t *testing.T) {
	cleanupTable(t, "tag")
	repo := NewTagRepo(testDB)

	// Создаем несколько тегов
	tag1, err := repo.Create(models.Tag{Name: "Tag1", HexColor: "#111111"})
	require.NoError(t, err)
	tag2, err := repo.Create(models.Tag{Name: "Tag2", HexColor: "#222222"})
	require.NoError(t, err)
	tag3, err := repo.Create(models.Tag{Name: "Tag3", HexColor: "#333333"})
	require.NoError(t, err)

	tests := []struct {
		name        string
		ids         []int64
		wantDeleted int
	}{
		{
			name:        "удаление нескольких тегов",
			ids:         []int64{tag1.ID, tag2.ID},
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

	// Проверяем, что tag3 все еще существует
	got, err := repo.Get(tag3.ID)
	require.NoError(t, err)
	assert.Equal(t, tag3.Name, got.Name)
}

func TestTagRepo_List(t *testing.T) {
	cleanupTable(t, "tag")
	repo := NewTagRepo(testDB)

	// Создаем тестовые данные
	tags := []models.Tag{
		{Name: "Alpha", HexColor: "#AAA111"},
		{Name: "Beta", HexColor: "#BBB222"},
		{Name: "Gamma", HexColor: "#CCC333"},
		{Name: "Delta", HexColor: "#DDD444"},
		{Name: "Epsilon", HexColor: "#EEE555"},
	}

	for _, tag := range tags {
		_, err := repo.Create(tag)
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
		{
			name: "сортировка по имени ASC",
			req: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 10,
				Sort: []entityreqdecorator.SortBy{
					{Field: "name", Order: "ASC"},
				},
			},
			wantTotal:   5,
			wantContent: 5,
		},
		{
			name: "сортировка по имени DESC",
			req: entityreqdecorator.PagebleRq{
				Page: 1,
				Size: 10,
				Sort: []entityreqdecorator.SortBy{
					{Field: "name", Order: "DESC"},
				},
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

func TestTagRepo_List_WithFilter(t *testing.T) {
	cleanupTable(t, "tag")
	repo := NewTagRepo(testDB)

	// Создаем тестовые данные
	_, err := repo.Create(models.Tag{Name: "Backend", HexColor: "#111111"})
	require.NoError(t, err)
	_, err = repo.Create(models.Tag{Name: "Frontend", HexColor: "#222222"})
	require.NoError(t, err)

	// Фильтрация по имени
	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
		Filter: map[string]entityreqdecorator.SQLGenerator{
			"name": &entityreqdecorator.PredicateLike{
				Predicate: entityreqdecorator.Predicate{Value: "Backend"},
			},
		},
	}

	result, err := repo.List(req)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Content, 1)
	assert.Equal(t, "Backend", result.Content[0].Name)
}

func TestTagRepo_List_Sorting(t *testing.T) {
	cleanupTable(t, "tag")
	repo := NewTagRepo(testDB)

	// Создаем тестовые данные в определенном порядке
	_, err := repo.Create(models.Tag{Name: "Charlie", HexColor: "#333333"})
	require.NoError(t, err)
	_, err = repo.Create(models.Tag{Name: "Alpha", HexColor: "#111111"})
	require.NoError(t, err)
	_, err = repo.Create(models.Tag{Name: "Bravo", HexColor: "#222222"})
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
	assert.Equal(t, "Alpha", result.Content[0].Name)
	assert.Equal(t, "Bravo", result.Content[1].Name)
	assert.Equal(t, "Charlie", result.Content[2].Name)

	// Сортировка по имени DESC
	req.Sort[0].Order = "DESC"
	result, err = repo.List(req)
	require.NoError(t, err)
	require.Len(t, result.Content, 3)
	assert.Equal(t, "Charlie", result.Content[0].Name)
	assert.Equal(t, "Bravo", result.Content[1].Name)
	assert.Equal(t, "Alpha", result.Content[2].Name)
}
