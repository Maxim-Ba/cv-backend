package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)




type TechnologyRepo struct {
	db *sql.DB
}

func NewTechnologyRepo(db *sql.DB) *TechnologyRepo {
	return &TechnologyRepo{
		db: db,
	}
}
// DeleteList удаляет список технологий по ID
func (t *TechnologyRepo) DeleteList(ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "DELETE FROM technology WHERE id = ANY($1) RETURNING id"
	rows, err := t.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to delete technology list: %w", err)
	}
	defer rows.Close()

	var deletedIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan deleted technology ID: %w", err)
		}
		deletedIDs = append(deletedIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return deletedIDs, nil
}

// Delete удаляет одну технологию по ID
func (t *TechnologyRepo) Delete(id int64) (int64, error) {
	query := "DELETE FROM technology WHERE id = $1"
	result, err := t.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete technology: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("technology with id %d not found", id)
	}

	return id, nil
}

// Get получает одну технологию по ID
func (t *TechnologyRepo) Get(id int64) (models.Technology, error) {
	query := "SELECT id, title, description, logo_url FROM technology WHERE id = $1"
	
	var technology models.Technology
	err := t.db.QueryRow(query, id).Scan(
		&technology.ID,
		&technology.Title,
		&technology.Description,
		&technology.LogoUrl,
	)
	
	if err == sql.ErrNoRows {
		return models.Technology{}, fmt.Errorf("technology with id %d not found", id)
	}
	if err != nil {
		return models.Technology{}, fmt.Errorf("failed to get technology: %w", err)
	}

	return technology, nil
}
func (t *TechnologyRepo) List(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Technology], error) {
	baseQuery := "SELECT id, title, description, logo_url FROM technology"

	queryParams := entityreqdecorator.BuildListQuery(
		req, baseQuery, t.isValidField,
	)

	var total int
	err := t.db.QueryRow(queryParams.CountQuery, queryParams.CountParams...).Scan(&total)
	if err != nil {
		return entityreqdecorator.PagebleRs[models.Technology]{}, fmt.Errorf("failed to count technologies: %w", err)
	}

	rows, err := t.db.Query(queryParams.SelectQuery, queryParams.SelectParams...)
	if err != nil {
		return entityreqdecorator.PagebleRs[models.Technology]{}, fmt.Errorf("failed to query technologies: %w", err)
	}
	defer rows.Close()

	var technologies []models.Technology
	for rows.Next() {
		var technology models.Technology
		err := rows.Scan(&technology.ID, &technology.Title, &technology.Description,  &technology.LogoUrl)
		if err != nil {
			return entityreqdecorator.PagebleRs[models.Technology]{}, fmt.Errorf("failed to scan technology: %w", err)
		}
		technologies = append(technologies, technology)
	}

	if err = rows.Err(); err != nil {
		return entityreqdecorator.PagebleRs[models.Technology]{}, fmt.Errorf("rows error: %w", err)
	}
	return entityreqdecorator.PagebleRs[models.Technology]{
		Total:   total,
		Content: technologies,
		Page:    req.Page,
		Size:    req.Size,
		Sort:    req.Sort,
	}, nil
}
// Create создает новую технологию
func (t *TechnologyRepo) Create(technology models.Technology) (models.Technology, error) {
	query := `
		INSERT INTO technology (title, description, logo_url)
		VALUES ($1, $2, $3)
		RETURNING id, title, description, logo_url
	`

	var created models.Technology
	err := t.db.QueryRow(
		query,
		technology.Title,
		technology.Description,
		technology.LogoUrl,
	).Scan(
		&created.ID,
		&created.Title,
		&created.Description,
		&created.LogoUrl,
	)

	if err != nil {
		return models.Technology{}, fmt.Errorf("failed to create technology: %w", err)
	}

	return created, nil
}

// Update обновляет существующую технологию
func (t *TechnologyRepo) Update(technology models.Technology) (models.Technology, error) {
	query := `
		UPDATE technology
		SET title = $2, description = $3, logo_url = $4
		WHERE id = $1
		RETURNING id, title, description, logo_url
	`

	var updated models.Technology
	err := t.db.QueryRow(
		query,
		technology.ID,
		technology.Title,
		technology.Description,
		technology.LogoUrl,
	).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Description,
		&updated.LogoUrl,
	)

	if err == sql.ErrNoRows {
		return models.Technology{}, fmt.Errorf("technology with id %d not found", technology.ID)
	}
	if err != nil {
		return models.Technology{}, fmt.Errorf("failed to update technology: %w", err)
	}

	return updated, nil
}

// isValidField проверяет, является ли поле валидным для сортировки и фильтрации
func (t *TechnologyRepo) isValidField(field string) bool {
	validFields := map[string]bool{
		"id":          true,
		"title":       true,
		"description": true,
		"logo_url":    true,
	}
	return validFields[field]
}
