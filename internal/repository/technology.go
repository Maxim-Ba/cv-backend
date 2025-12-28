package repository

import (
	"database/sql"
	"fmt"

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
func (t *TechnologyRepo) DeleteList([]int64) ([]int64, error) {
	return nil, nil
}
func (t *TechnologyRepo) Delete(id int64) (int64, error) {
	return 0, nil
}
func (t *TechnologyRepo) Get(id int64) (models.Technology, error) {
	return models.Technology{}, nil
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
func (t *TechnologyRepo) Update() (models.Technology, error) {
	return models.Technology{}, nil
}
func (t *TechnologyRepo) Create() (models.Technology, error) {
	return models.Technology{}, nil
}

func (t *TechnologyRepo) isValidField(field string) bool {
	validFields := map[string]bool{
		"id":        true,
		"name":      true,
		"hex_color": true,
	}
	return validFields[field]
}
