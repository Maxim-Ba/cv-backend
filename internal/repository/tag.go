package repository

import (
	"database/sql"
	"fmt"

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
func (t *TagRepo) DeleteList([]int64) ([]int64, error) {
	return nil, nil
}
func (t *TagRepo) Delete(id int64) (int64, error) {
	return 0, nil
}
func (t *TagRepo) Get(id int64) (models.Tag, error) {
	return models.Tag{}, nil
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
func (t *TagRepo) Update() (models.Tag, error) {
	return models.Tag{}, nil
}
func (t *TagRepo) Create() (models.Tag, error) {
	return models.Tag{}, nil
}

func (t *TagRepo) isValidField(field string) bool {
	validFields := map[string]bool{
		"id":        true,
		"name":      true,
		"hex_color": true,
	}
	return validFields[field]
}
