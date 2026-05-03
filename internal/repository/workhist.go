package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/Maxim-Ba/cv-backend/internal/models/dto"
	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

// WorkHistoryRepo репозиторий для работы с таблицей work_history
type WorkHistoryRepo struct {
	db *sql.DB
}

// NewWorkHistoryRepo создает новый экземпляр репозитория истории работы
func NewWorkHistoryRepo(db *sql.DB) *WorkHistoryRepo {
	return &WorkHistoryRepo{
		db: db,
	}
}

// DeleteList удаляет список записей истории работы по ID
func (w *WorkHistoryRepo) DeleteList(ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "DELETE FROM work_history WHERE id = ANY($1) RETURNING id"
	rows, err := w.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to delete work history list: %w", err)
	}
	defer rows.Close()

	var deletedIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan deleted work history ID: %w", err)
		}
		deletedIDs = append(deletedIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return deletedIDs, nil
}

// Delete удаляет одну запись истории работы по ID
func (w *WorkHistoryRepo) Delete(id int64) (int64, error) {
	query := "DELETE FROM work_history WHERE id = $1"
	result, err := w.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete work history: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("work history with id %d not found", id)
	}

	return id, nil
}

// Get получает одну запись истории работы по ID
func (w *WorkHistoryRepo) Get(id int64) (models.WorkHistory, error) {
	query := `
		SELECT id, name, about, logo_url, period_start, period_end, what_i_did, projects
		FROM work_history
		WHERE id = $1
	`
	
	var workHistory models.WorkHistory
	err := w.db.QueryRow(query, id).Scan(
		&workHistory.ID,
		&workHistory.Name,
		&workHistory.About,
		&workHistory.LogoUrl,
		&workHistory.PeriodStart,
		&workHistory.PeriodEnd,
		pq.Array(&workHistory.WhatIDid),
		pq.Array(&workHistory.Projects),
	)
	
	if err == sql.ErrNoRows {
		return models.WorkHistory{}, fmt.Errorf("work history with id %d not found", id)
	}
	if err != nil {
		return models.WorkHistory{}, fmt.Errorf("failed to get work history: %w", err)
	}

	return workHistory, nil
}

// List получает список записей истории работы с пагинацией, сортировкой и фильтрацией
func (w *WorkHistoryRepo) List(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.WorkHistory], error) {
	baseQuery := `
		SELECT id, name, about, logo_url, period_start, period_end, what_i_did, projects
		FROM work_history
	`

	queryParams := entityreqdecorator.BuildListQuery(
		req, baseQuery, w.isValidField,
	)

	// Получаем общее количество записей
	var total int
	err := w.db.QueryRow(queryParams.CountQuery, queryParams.CountParams...).Scan(&total)
	if err != nil {
		return entityreqdecorator.PagebleRs[models.WorkHistory]{}, fmt.Errorf("failed to count work histories: %w", err)
	}

	// Получаем записи с учетом пагинации
	rows, err := w.db.Query(queryParams.SelectQuery, queryParams.SelectParams...)
	if err != nil {
		return entityreqdecorator.PagebleRs[models.WorkHistory]{}, fmt.Errorf("failed to query work histories: %w", err)
	}
	defer rows.Close()

	var workHistories []models.WorkHistory
	for rows.Next() {
		var workHistory models.WorkHistory
		err := rows.Scan(
			&workHistory.ID,
			&workHistory.Name,
			&workHistory.About,
			&workHistory.LogoUrl,
			&workHistory.PeriodStart,
			&workHistory.PeriodEnd,
			pq.Array(&workHistory.WhatIDid),
			pq.Array(&workHistory.Projects),
		)
		if err != nil {
			return entityreqdecorator.PagebleRs[models.WorkHistory]{}, fmt.Errorf("failed to scan work history: %w", err)
		}
		workHistories = append(workHistories, workHistory)
	}

	if err = rows.Err(); err != nil {
		return entityreqdecorator.PagebleRs[models.WorkHistory]{}, fmt.Errorf("rows error: %w", err)
	}

	return entityreqdecorator.PagebleRs[models.WorkHistory]{
		Total:   total,
		Content: workHistories,
		Page:    req.Page,
		Size:    req.Size,
		Sort:    req.Sort,
	}, nil
}

// Create создает новую запись истории работы
func (w *WorkHistoryRepo) Create(workHistory models.WorkHistory) (models.WorkHistory, error) {
	query := `
		INSERT INTO work_history (name, about, logo_url, period_start, period_end, what_i_did, projects)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, about, logo_url, period_start, period_end, what_i_did, projects
	`

	var created models.WorkHistory
	err := w.db.QueryRow(
		query,
		workHistory.Name,
		workHistory.About,
		workHistory.LogoUrl,
		workHistory.PeriodStart,
		workHistory.PeriodEnd,
		pq.Array(workHistory.WhatIDid),
		pq.Array(workHistory.Projects),
	).Scan(
		&created.ID,
		&created.Name,
		&created.About,
		&created.LogoUrl,
		&created.PeriodStart,
		&created.PeriodEnd,
		pq.Array(&created.WhatIDid),
		pq.Array(&created.Projects),
	)

	if err != nil {
		return models.WorkHistory{}, fmt.Errorf("failed to create work history: %w", err)
	}

	return created, nil
}

// Update обновляет существующую запись истории работы
func (w *WorkHistoryRepo) Update(workHistory models.WorkHistory) (models.WorkHistory, error) {
	query := `
		UPDATE work_history
		SET name = $2, about = $3, logo_url = $4, period_start = $5, 
		    period_end = $6, what_i_did = $7, projects = $8
		WHERE id = $1
		RETURNING id, name, about, logo_url, period_start, period_end, what_i_did, projects
	`

	var updated models.WorkHistory
	err := w.db.QueryRow(
		query,
		workHistory.ID,
		workHistory.Name,
		workHistory.About,
		workHistory.LogoUrl,
		workHistory.PeriodStart,
		workHistory.PeriodEnd,
		pq.Array(workHistory.WhatIDid),
		pq.Array(workHistory.Projects),
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.About,
		&updated.LogoUrl,
		&updated.PeriodStart,
		&updated.PeriodEnd,
		pq.Array(&updated.WhatIDid),
		pq.Array(&updated.Projects),
	)

	if err == sql.ErrNoRows {
		return models.WorkHistory{}, fmt.Errorf("work history with id %d not found", workHistory.ID)
	}
	if err != nil {
		return models.WorkHistory{}, fmt.Errorf("failed to update work history: %w", err)
	}

	return updated, nil
}

// isValidField проверяет, является ли поле валидным для сортировки и фильтрации
func (w *WorkHistoryRepo) isValidField(field string) bool {
	validFields := map[string]bool{
		"id":           true,
		"name":         true,
		"about":        true,
		"period_start": true,
		"period_end":   true,
	}
	return validFields[field]
}

// nullTimeToString конвертирует sql.NullTime в *string формата "2006-01-02"
func nullTimeToString(t sql.NullTime) *string {
	if !t.Valid {
		return nil
	}
	s := t.Time.Format(time.DateOnly)
	return &s
}

// GetWithTechnologies получает историю работы по ID вместе с технологиями и их тегами
func (w *WorkHistoryRepo) GetWithTechnologies(id int64) (dto.WorkHistoryWithTechnologiesDTO, error) {
	query := `
		SELECT wh.id, wh.name, wh.about, wh.logo_url,
		       wh.period_start, wh.period_end,
		       wh.what_i_did, wh.projects,
		       t.id, t.title, t.description, t.logo_url,
		       tg.id, tg.name, tg.hex_color
		FROM work_history wh
		LEFT JOIN work_history_technology wht ON wht.work_history_id = wh.id
		LEFT JOIN technology t ON t.id = wht.technology_id
		LEFT JOIN technologies_tag tt ON tt.technology_id = t.id
		LEFT JOIN tag tg ON tg.id = tt.tag_id
		WHERE wh.id = $1
		ORDER BY t.id
	`

	rows, err := w.db.Query(query, id)
	if err != nil {
		return dto.WorkHistoryWithTechnologiesDTO{}, fmt.Errorf("failed to get work history with technologies: %w", err)
	}
	defer rows.Close()

	var result *dto.WorkHistoryWithTechnologiesDTO
	techMap := make(map[int64]*dto.TechnologyWithTagsDTO)

	for rows.Next() {
		var (
			whID        int64
			name        string
			about       string
			logoUrlBytes []byte
			periodStart sql.NullTime
			periodEnd   sql.NullTime
			whatIDid    []string
			projects    []string
			techID      sql.NullInt64
			techTitle   sql.NullString
			techDesc    sql.NullString
			techLogo    sql.NullString
			tagID       sql.NullInt64
			tagName     sql.NullString
			tagHexColor sql.NullString
		)
		if err := rows.Scan(
			&whID, &name, &about, &logoUrlBytes,
			&periodStart, &periodEnd,
			pq.Array(&whatIDid), pq.Array(&projects),
			&techID, &techTitle, &techDesc, &techLogo,
			&tagID, &tagName, &tagHexColor,
		); err != nil {
			return dto.WorkHistoryWithTechnologiesDTO{}, fmt.Errorf("failed to scan work history row: %w", err)
		}
		if result == nil {
			result = &dto.WorkHistoryWithTechnologiesDTO{
				WorkHistoryDTO: dto.WorkHistoryDTO{
					ID:          whID,
					Name:        name,
					About:       about,
					LogoUrl:     string(logoUrlBytes),
					PeriodStart: nullTimeToString(periodStart),
					PeriodEnd:   nullTimeToString(periodEnd),
					WhatIDid:    whatIDid,
					Projects:    projects,
				},
				Technologies: []dto.TechnologyWithTagsDTO{},
			}
		}
		if techID.Valid {
			if _, exists := techMap[techID.Int64]; !exists {
				tech := dto.TechnologyWithTagsDTO{
					TechnologyDTO: dto.TechnologyDTO{
						ID:    techID.Int64,
						Title: techTitle.String,
					},
					Tags: []dto.TagDTO{},
				}
				if techDesc.Valid {
					tech.Description = &techDesc.String
				}
				if techLogo.Valid {
					tech.LogoUrl = &techLogo.String
				}
				techMap[techID.Int64] = &tech
				result.Technologies = append(result.Technologies, tech)
			}
			if tagID.Valid {
				for i := range result.Technologies {
					if result.Technologies[i].ID == techID.Int64 {
						result.Technologies[i].Tags = append(result.Technologies[i].Tags, dto.TagDTO{
							ID:       tagID.Int64,
							Name:     tagName.String,
							HexColor: tagHexColor.String,
						})
						break
					}
				}
			}
		}
	}
	if err := rows.Err(); err != nil {
		return dto.WorkHistoryWithTechnologiesDTO{}, fmt.Errorf("rows error: %w", err)
	}
	if result == nil {
		return dto.WorkHistoryWithTechnologiesDTO{}, fmt.Errorf("work history with id %d not found", id)
	}
	return *result, nil
}

// ListWithTechnologies получает список истории работы с технологиями и их тегами
func (w *WorkHistoryRepo) ListWithTechnologies(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO], error) {
	var total int
	if err := w.db.QueryRow("SELECT COUNT(*) FROM work_history").Scan(&total); err != nil {
		return entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO]{}, fmt.Errorf("failed to count work histories: %w", err)
	}

	offset := 0
	limit := total
	if req.Size > 0 {
		limit = req.Size
		offset = req.Page * req.Size
	}

	// Получаем ID истории работы с пагинацией
	idRows, err := w.db.Query("SELECT id FROM work_history ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO]{}, fmt.Errorf("failed to query work history ids: %w", err)
	}
	defer idRows.Close()

	var ids []int64
	for idRows.Next() {
		var id int64
		if err := idRows.Scan(&id); err != nil {
			return entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO]{}, fmt.Errorf("failed to scan work history id: %w", err)
		}
		ids = append(ids, id)
	}
	if err := idRows.Err(); err != nil {
		return entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO]{}, fmt.Errorf("ids rows error: %w", err)
	}

	if len(ids) == 0 {
		return entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO]{
			Total:   total,
			Content: []dto.WorkHistoryWithTechnologiesDTO{},
			Page:    req.Page,
			Size:    req.Size,
		}, nil
	}

	query := `
		SELECT wh.id, wh.name, wh.about, wh.logo_url,
		       wh.period_start, wh.period_end,
		       wh.what_i_did, wh.projects,
		       t.id, t.title, t.description, t.logo_url,
		       tg.id, tg.name, tg.hex_color
		FROM work_history wh
		LEFT JOIN work_history_technology wht ON wht.work_history_id = wh.id
		LEFT JOIN technology t ON t.id = wht.technology_id
		LEFT JOIN technologies_tag tt ON tt.technology_id = t.id
		LEFT JOIN tag tg ON tg.id = tt.tag_id
		WHERE wh.id = ANY($1)
		ORDER BY wh.id, t.id
	`

	rows, err := w.db.Query(query, pq.Array(ids))
	if err != nil {
		return entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO]{}, fmt.Errorf("failed to query work histories with technologies: %w", err)
	}
	defer rows.Close()

	whMap := make(map[int64]*dto.WorkHistoryWithTechnologiesDTO)
	// techTagMap[whID][techID] = tech index in Technologies slice
	type techKey struct{ whID, techID int64 }
	techTagMap := make(map[techKey]int)

	for rows.Next() {
		var (
			whID         int64
			name         string
			about        string
			logoUrlBytes []byte
			periodStart  sql.NullTime
			periodEnd    sql.NullTime
			whatIDid     []string
			projects     []string
			techID       sql.NullInt64
			techTitle    sql.NullString
			techDesc     sql.NullString
			techLogo     sql.NullString
			tagID        sql.NullInt64
			tagName      sql.NullString
			tagHexColor  sql.NullString
		)
		if err := rows.Scan(
			&whID, &name, &about, &logoUrlBytes,
			&periodStart, &periodEnd,
			pq.Array(&whatIDid), pq.Array(&projects),
			&techID, &techTitle, &techDesc, &techLogo,
			&tagID, &tagName, &tagHexColor,
		); err != nil {
			return entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO]{}, fmt.Errorf("failed to scan row: %w", err)
		}

		if _, exists := whMap[whID]; !exists {
			whMap[whID] = &dto.WorkHistoryWithTechnologiesDTO{
				WorkHistoryDTO: dto.WorkHistoryDTO{
					ID:          whID,
					Name:        name,
					About:       about,
					LogoUrl:     string(logoUrlBytes),
					PeriodStart: nullTimeToString(periodStart),
					PeriodEnd:   nullTimeToString(periodEnd),
					WhatIDid:    whatIDid,
					Projects:    projects,
				},
				Technologies: []dto.TechnologyWithTagsDTO{},
			}
		}

		if techID.Valid {
			key := techKey{whID: whID, techID: techID.Int64}
			if _, exists := techTagMap[key]; !exists {
				tech := dto.TechnologyWithTagsDTO{
					TechnologyDTO: dto.TechnologyDTO{
						ID:    techID.Int64,
						Title: techTitle.String,
					},
					Tags: []dto.TagDTO{},
				}
				if techDesc.Valid {
					tech.Description = &techDesc.String
				}
				if techLogo.Valid {
					tech.LogoUrl = &techLogo.String
				}
				whMap[whID].Technologies = append(whMap[whID].Technologies, tech)
				techTagMap[key] = len(whMap[whID].Technologies) - 1
			}

			if tagID.Valid {
				idx := techTagMap[key]
				whMap[whID].Technologies[idx].Tags = append(whMap[whID].Technologies[idx].Tags, dto.TagDTO{
					ID:       tagID.Int64,
					Name:     tagName.String,
					HexColor: tagHexColor.String,
				})
			}
		}
	}
	if err := rows.Err(); err != nil {
		return entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO]{}, fmt.Errorf("rows error: %w", err)
	}

	result := make([]dto.WorkHistoryWithTechnologiesDTO, 0, len(ids))
	for _, id := range ids {
		if wh, ok := whMap[id]; ok {
			result = append(result, *wh)
		}
	}

	return entityreqdecorator.PagebleRs[dto.WorkHistoryWithTechnologiesDTO]{
		Total:   total,
		Content: result,
		Page:    req.Page,
		Size:    req.Size,
		Sort:    req.Sort,
	}, nil
}
