package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

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
