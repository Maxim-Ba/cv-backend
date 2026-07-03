package repository

import (
	"testing"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEducationRepo_Create(t *testing.T) {
	cleanupTable(t, "education")
	repo := NewEducationRepo(testDB)

	tests := []struct {
		name      string
		education models.Education
		wantErr   bool
	}{
		{
			name: "успешное создание записи образования",
			education: models.Education{
				Name:         newNullableLocalizedText("Основы Go"),
				Year:         2023,
				Course:       newLocalizedText("Backend Development"),
				Organization: newLocalizedText("Yandex Practicum"),
			},
			wantErr: false,
		},
		{
			name: "успешное создание записи без имени",
			education: models.Education{
				Name:         i18n.NullableLocalizedText{},
				Year:         2022,
				Course:       newLocalizedText("Python"),
				Organization: newLocalizedText("Coursera"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := repo.Create(tt.education)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotZero(t, created.ID)
			assert.Equal(t, tt.education.Name, created.Name)
			assert.Equal(t, tt.education.Year, created.Year)
			assert.Equal(t, tt.education.Course, created.Course)
			assert.Equal(t, tt.education.Organization, created.Organization)
		})
	}
}

func TestEducationRepo_Get(t *testing.T) {
	cleanupTable(t, "education")
	repo := NewEducationRepo(testDB)

	created, err := repo.Create(models.Education{
		Name:         newNullableLocalizedText("Test Course"),
		Year:         2024,
		Course:       newLocalizedText("Testing"),
		Organization: newLocalizedText("Test Org"),
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
			assert.Equal(t, created.Year, got.Year)
			assert.Equal(t, created.Course, got.Course)
			assert.Equal(t, created.Organization, got.Organization)
		})
	}
}

func TestEducationRepo_Update(t *testing.T) {
	cleanupTable(t, "education")
	repo := NewEducationRepo(testDB)

	created, err := repo.Create(models.Education{
		Name:         newNullableLocalizedText("Original"),
		Year:         2020,
		Course:       newLocalizedText("Original Course"),
		Organization: newLocalizedText("Original Org"),
	})
	require.NoError(t, err)

	tests := []struct {
		name      string
		education models.Education
		wantErr   bool
	}{
		{
			name: "успешное обновление записи",
			education: models.Education{
				ID:           created.ID,
				Name:         newNullableLocalizedText("Updated"),
				Year:         2024,
				Course:       newLocalizedText("Updated Course"),
				Organization: newLocalizedText("Updated Org"),
			},
			wantErr: false,
		},
		{
			name: "обновление несуществующей записи",
			education: models.Education{
				ID:           99999,
				Name:         newNullableLocalizedText("NonExistent"),
				Year:         2000,
				Course:       newLocalizedText("No Course"),
				Organization: newLocalizedText("No Org"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := repo.Update(tt.education)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.education.ID, updated.ID)
			assert.Equal(t, tt.education.Name, updated.Name)
			assert.Equal(t, tt.education.Year, updated.Year)
			assert.Equal(t, tt.education.Course, updated.Course)
			assert.Equal(t, tt.education.Organization, updated.Organization)
		})
	}
}

func TestEducationRepo_Delete(t *testing.T) {
	cleanupTable(t, "education")
	repo := NewEducationRepo(testDB)

	created, err := repo.Create(models.Education{
		Name:         newNullableLocalizedText("ToDelete"),
		Year:         2021,
		Course:       newLocalizedText("Delete Course"),
		Organization: newLocalizedText("Delete Org"),
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

			_, err = repo.Get(tt.id)
			require.Error(t, err)
		})
	}
}

func TestEducationRepo_DeleteList(t *testing.T) {
	cleanupTable(t, "education")
	repo := NewEducationRepo(testDB)

	edu1, err := repo.Create(models.Education{
		Name:         newNullableLocalizedText("Edu1"),
		Year:         2020,
		Course:       newLocalizedText("Course1"),
		Organization: newLocalizedText("Org1"),
	})
	require.NoError(t, err)

	edu2, err := repo.Create(models.Education{
		Name:         newNullableLocalizedText("Edu2"),
		Year:         2021,
		Course:       newLocalizedText("Course2"),
		Organization: newLocalizedText("Org2"),
	})
	require.NoError(t, err)

	edu3, err := repo.Create(models.Education{
		Name:         newNullableLocalizedText("Edu3"),
		Year:         2022,
		Course:       newLocalizedText("Course3"),
		Organization: newLocalizedText("Org3"),
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		ids         []int64
		wantDeleted int
	}{
		{
			name:        "удаление нескольких записей",
			ids:         []int64{edu1.ID, edu2.ID},
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

	got, err := repo.Get(edu3.ID)
	require.NoError(t, err)
	assert.Equal(t, "Edu3", got.Name.Get(i18n.LocaleRU))
}

func TestEducationRepo_List(t *testing.T) {
	cleanupTable(t, "education")
	repo := NewEducationRepo(testDB)

	educations := []models.Education{
		{Name: newNullableLocalizedText("Course A"), Year: 2020, Course: newLocalizedText("A"), Organization: newLocalizedText("Org A")},
		{Name: newNullableLocalizedText("Course B"), Year: 2021, Course: newLocalizedText("B"), Organization: newLocalizedText("Org B")},
		{Name: newNullableLocalizedText("Course C"), Year: 2022, Course: newLocalizedText("C"), Organization: newLocalizedText("Org C")},
		{Name: newNullableLocalizedText("Course D"), Year: 2023, Course: newLocalizedText("D"), Organization: newLocalizedText("Org D")},
		{Name: newNullableLocalizedText("Course E"), Year: 2024, Course: newLocalizedText("E"), Organization: newLocalizedText("Org E")},
	}

	for _, edu := range educations {
		_, err := repo.Create(edu)
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

func TestEducationRepo_List_WithFilter(t *testing.T) {
	cleanupTable(t, "education")
	repo := NewEducationRepo(testDB)

	_, err := repo.Create(models.Education{
		Name:         newNullableLocalizedText("Go Course"),
		Year:         2023,
		Course:       newLocalizedText("Go Programming"),
		Organization: newLocalizedText("Yandex"),
	})
	require.NoError(t, err)

	_, err = repo.Create(models.Education{
		Name:         newNullableLocalizedText("Python Course"),
		Year:         2022,
		Course:       newLocalizedText("Python Programming"),
		Organization: newLocalizedText("Coursera"),
	})
	require.NoError(t, err)

	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
		Filter: map[string]entityreqdecorator.SQLGenerator{
			"year": &entityreqdecorator.PredicateEQ{
				Predicate: entityreqdecorator.Predicate{Value: "2023"},
			},
		},
	}

	result, err := repo.List(req)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Content, 1)
	assert.Equal(t, int32(2023), result.Content[0].Year)
}

func TestEducationRepo_List_Sorting(t *testing.T) {
	cleanupTable(t, "education")
	repo := NewEducationRepo(testDB)

	_, err := repo.Create(models.Education{
		Name:         newNullableLocalizedText("2022 Course"),
		Year:         2022,
		Course:       newLocalizedText("Course"),
		Organization: newLocalizedText("Org"),
	})
	require.NoError(t, err)

	_, err = repo.Create(models.Education{
		Name:         newNullableLocalizedText("2024 Course"),
		Year:         2024,
		Course:       newLocalizedText("Course"),
		Organization: newLocalizedText("Org"),
	})
	require.NoError(t, err)

	_, err = repo.Create(models.Education{
		Name:         newNullableLocalizedText("2023 Course"),
		Year:         2023,
		Course:       newLocalizedText("Course"),
		Organization: newLocalizedText("Org"),
	})
	require.NoError(t, err)

	req := entityreqdecorator.PagebleRq{
		Page: 1,
		Size: 10,
		Sort: []entityreqdecorator.SortBy{
			{Field: "year", Order: "ASC"},
		},
	}

	result, err := repo.List(req)
	require.NoError(t, err)
	require.Len(t, result.Content, 3)
	assert.Equal(t, int32(2022), result.Content[0].Year)
	assert.Equal(t, int32(2023), result.Content[1].Year)
	assert.Equal(t, int32(2024), result.Content[2].Year)

	req.Sort[0].Order = "DESC"
	result, err = repo.List(req)
	require.NoError(t, err)
	require.Len(t, result.Content, 3)
	assert.Equal(t, int32(2024), result.Content[0].Year)
	assert.Equal(t, int32(2023), result.Content[1].Year)
	assert.Equal(t, int32(2022), result.Content[2].Year)
}
