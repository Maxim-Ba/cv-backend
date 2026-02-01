package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

type TagRepo struct {
	db *sql.DB
}

func NewTagRepo(db *sql.DB) *TagRepo {
	return &TagRepo{
		db: db,
	}
}
// DeleteList удаляет список тегов по ID
func (t *TagRepo) DeleteList(ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "DELETE FROM tag WHERE id = ANY($1) RETURNING id"
	rows, err := t.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to delete tag list: %w", err)
	}
	defer rows.Close()

	var deletedIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan deleted tag ID: %w", err)
		}
		deletedIDs = append(deletedIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return deletedIDs, nil
}

// Delete удаляет один тег по ID
func (t *TagRepo) Delete(id int64) (int64, error) {
	query := "DELETE FROM tag WHERE id = $1"
	result, err := t.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("tag with id %d not found", id)
	}

	return id, nil
}

// Get получает один тег по ID
func (t *TagRepo) Get(id int64) (models.Tag, error) {
	query := "SELECT id, name, hex_color FROM tag WHERE id = $1"
	
	var tag models.Tag
	err := t.db.QueryRow(query, id).Scan(&tag.ID, &tag.Name, &tag.HexColor)
	
	if err == sql.ErrNoRows {
		return models.Tag{}, fmt.Errorf("tag with id %d not found", id)
	}
	if err != nil {
		return models.Tag{}, fmt.Errorf("failed to get tag: %w", err)
	}

	return tag, nil
}
func (t *TagRepo) List(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Tag], error) {
	baseQuery := "SELECT id, name, hex_color FROM tag"

	queryParams := entityreqdecorator.BuildListQuery(
		req, baseQuery, t.isValidField,
	)

	var total int
	err := t.db.QueryRow(queryParams.CountQuery, queryParams.CountParams...).Scan(&total)
	if err != nil {
		return entityreqdecorator.PagebleRs[models.Tag]{}, fmt.Errorf("failed to count tags: %w", err)
	}

	rows, err := t.db.Query(queryParams.SelectQuery, queryParams.SelectParams...)
	if err != nil {
		return entityreqdecorator.PagebleRs[models.Tag]{}, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.HexColor)
		if err != nil {
			return entityreqdecorator.PagebleRs[models.Tag]{}, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return entityreqdecorator.PagebleRs[models.Tag]{}, fmt.Errorf("rows error: %w", err)
	}
	return entityreqdecorator.PagebleRs[models.Tag]{
		Total:   total,
		Content: tags,
		Page:    req.Page,
		Size:    req.Size,
		Sort:    req.Sort,
	}, nil
}
// Create создает новый тег
func (t *TagRepo) Create(tag models.Tag) (models.Tag, error) {
	query := `
		INSERT INTO tag (name, hex_color)
		VALUES ($1, $2)
		RETURNING id, name, hex_color
	`

	var created models.Tag
	err := t.db.QueryRow(query, tag.Name, tag.HexColor).Scan(
		&created.ID,
		&created.Name,
		&created.HexColor,
	)

	if err != nil {
		return models.Tag{}, fmt.Errorf("failed to create tag: %w", err)
	}

	return created, nil
}

// Update обновляет существующий тег
func (t *TagRepo) Update(tag models.Tag) (models.Tag, error) {
	query := `
		UPDATE tag
		SET name = $2, hex_color = $3
		WHERE id = $1
		RETURNING id, name, hex_color
	`

	var updated models.Tag
	err := t.db.QueryRow(query, tag.ID, tag.Name, tag.HexColor).Scan(
		&updated.ID,
		&updated.Name,
		&updated.HexColor,
	)

	if err == sql.ErrNoRows {
		return models.Tag{}, fmt.Errorf("tag with id %d not found", tag.ID)
	}
	if err != nil {
		return models.Tag{}, fmt.Errorf("failed to update tag: %w", err)
	}

	return updated, nil
}

func (t *TagRepo) isValidField(field string) bool {
	validFields := map[string]bool{
		"id":        true,
		"name":      true,
		"hex_color": true,
	}
	return validFields[field]
}
