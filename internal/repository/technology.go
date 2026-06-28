package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/Maxim-Ba/cv-backend/internal/models/dto"
	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	"github.com/Maxim-Ba/cv-backend/pkg/apierrors"
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
		return 0, fmt.Errorf("technology with id %d: %w", id, apierrors.ErrNotFound)
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
		return models.Technology{}, fmt.Errorf("technology with id %d: %w", id, apierrors.ErrNotFound)
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
		err := rows.Scan(&technology.ID, &technology.Title, &technology.Description, &technology.LogoUrl)
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
		return models.Technology{}, fmt.Errorf("technology with id %d: %w", technology.ID, apierrors.ErrNotFound)
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

// technologyRowToDTO конвертирует поля строки запроса в TechnologyDTO
func technologyRowToDTO(id int64, title string, description sql.NullString, logoUrl sql.NullString) dto.TechnologyWithTagsDTO {
	tech := dto.TechnologyWithTagsDTO{
		TechnologyDTO: dto.TechnologyDTO{
			ID:    id,
			Title: title,
		},
		Tags: []dto.TagDTO{},
	}
	if description.Valid {
		tech.Description = &description.String
	}
	if logoUrl.Valid {
		tech.LogoUrl = &logoUrl.String
	}
	return tech
}

// GetWithTags получает технологию по ID вместе с её тегами
func (t *TechnologyRepo) GetWithTags(id int64) (dto.TechnologyWithTagsDTO, error) {
	query := `
		SELECT t.id, t.title, t.description, t.logo_url,
		       tg.id, tg.name, tg.hex_color
		FROM technology t
		LEFT JOIN technologies_tag tt ON tt.technology_id = t.id
		LEFT JOIN tag tg ON tg.id = tt.tag_id
		WHERE t.id = $1
	`

	rows, err := t.db.Query(query, id)
	if err != nil {
		return dto.TechnologyWithTagsDTO{}, fmt.Errorf("failed to get technology with tags: %w", err)
	}
	defer rows.Close()

	var result *dto.TechnologyWithTagsDTO
	for rows.Next() {
		var (
			techID      int64
			title       string
			description sql.NullString
			logoUrl     sql.NullString
			tagID       sql.NullInt64
			tagName     sql.NullString
			tagHexColor sql.NullString
		)
		if err := rows.Scan(&techID, &title, &description, &logoUrl, &tagID, &tagName, &tagHexColor); err != nil {
			return dto.TechnologyWithTagsDTO{}, fmt.Errorf("failed to scan technology with tags: %w", err)
		}
		if result == nil {
			t := technologyRowToDTO(techID, title, description, logoUrl)
			result = &t
		}
		if tagID.Valid {
			result.Tags = append(result.Tags, dto.TagDTO{
				ID:       tagID.Int64,
				Name:     tagName.String,
				HexColor: tagHexColor.String,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return dto.TechnologyWithTagsDTO{}, fmt.Errorf("rows error: %w", err)
	}
	if result == nil {
		return dto.TechnologyWithTagsDTO{}, fmt.Errorf("technology with id %d not found", id)
	}
	return *result, nil
}

// ListWithTags получает список технологий с тегами и пагинацией
func (t *TechnologyRepo) ListWithTags(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO], error) {
	var total int
	if err := t.db.QueryRow("SELECT COUNT(*) FROM technology").Scan(&total); err != nil {
		return entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO]{}, fmt.Errorf("failed to count technologies: %w", err)
	}

	offset := 0
	limit := total
	if req.Size > 0 {
		limit = req.Size
		offset = (req.Page - 1) * req.Size
	}

	idRows, err := t.db.Query("SELECT id FROM technology ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO]{}, fmt.Errorf("failed to query technology ids: %w", err)
	}
	defer idRows.Close()

	var ids []int64
	for idRows.Next() {
		var id int64
		if err := idRows.Scan(&id); err != nil {
			return entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO]{}, fmt.Errorf("failed to scan technology id: %w", err)
		}
		ids = append(ids, id)
	}
	if err := idRows.Err(); err != nil {
		return entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO]{}, fmt.Errorf("technology ids rows error: %w", err)
	}

	if len(ids) == 0 {
		return entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO]{
			Total:   total,
			Content: []dto.TechnologyWithTagsDTO{},
			Page:    req.Page,
			Size:    req.Size,
			Sort:    req.Sort,
		}, nil
	}

	selectQuery := `
		SELECT t.id, t.title, t.description, t.logo_url,
		       tg.id, tg.name, tg.hex_color
		FROM technology t
		LEFT JOIN technologies_tag tt ON tt.technology_id = t.id
		LEFT JOIN tag tg ON tg.id = tt.tag_id
		WHERE t.id = ANY($1)
		ORDER BY t.id
	`

	rows, err := t.db.Query(selectQuery, pq.Array(ids))
	if err != nil {
		return entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO]{}, fmt.Errorf("failed to query technologies with tags: %w", err)
	}
	defer rows.Close()

	techMap := make(map[int64]*dto.TechnologyWithTagsDTO)
	var order []int64

	for rows.Next() {
		var (
			techID      int64
			title       string
			description sql.NullString
			logoUrl     sql.NullString
			tagID       sql.NullInt64
			tagName     sql.NullString
			tagHexColor sql.NullString
		)
		if err := rows.Scan(&techID, &title, &description, &logoUrl, &tagID, &tagName, &tagHexColor); err != nil {
			return entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO]{}, fmt.Errorf("failed to scan technology row: %w", err)
		}
		if _, exists := techMap[techID]; !exists {
			tech := technologyRowToDTO(techID, title, description, logoUrl)
			techMap[techID] = &tech
			order = append(order, techID)
		}
		if tagID.Valid {
			techMap[techID].Tags = append(techMap[techID].Tags, dto.TagDTO{
				ID:       tagID.Int64,
				Name:     tagName.String,
				HexColor: tagHexColor.String,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO]{}, fmt.Errorf("rows error: %w", err)
	}

	technologies := make([]dto.TechnologyWithTagsDTO, 0, len(order))
	for _, id := range order {
		technologies = append(technologies, *techMap[id])
	}

	return entityreqdecorator.PagebleRs[dto.TechnologyWithTagsDTO]{
		Total:   total,
		Content: technologies,
		Page:    req.Page,
		Size:    req.Size,
		Sort:    req.Sort,
	}, nil
}

// SetTags заменяет набор тегов у технологии
func (t *TechnologyRepo) SetTags(technologyID int64, tagIDs []int64) error {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec("DELETE FROM technologies_tag WHERE technology_id = $1", technologyID); err != nil {
		return fmt.Errorf("failed to remove existing tags: %w", err)
	}

	for _, tagID := range tagIDs {
		if _, err = tx.Exec(
			"INSERT INTO technologies_tag (tag_id, technology_id) VALUES ($1, $2)",
			tagID,
			technologyID,
		); err != nil {
			return fmt.Errorf("failed to add tag %d to technology %d: %w", tagID, technologyID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tags transaction: %w", err)
	}

	return nil
}
