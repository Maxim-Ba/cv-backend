package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

// EducationRepo репозиторий для работы с таблицей education
type EducationRepo struct {
	db *sql.DB
}

// NewEducationRepo создает новый экземпляр репозитория образования
func NewEducationRepo(db *sql.DB) *EducationRepo {
	return &EducationRepo{
		db: db,
	}
}

// DeleteList удаляет список записей образования по ID
func (e *EducationRepo) DeleteList(ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "DELETE FROM education WHERE id = ANY($1) RETURNING id"
	rows, err := e.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to delete education list: %w", err)
	}
	defer rows.Close()

	var deletedIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan deleted education ID: %w", err)
		}
		deletedIDs = append(deletedIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return deletedIDs, nil
}

// Delete удаляет одну запись образования по ID
func (e *EducationRepo) Delete(id int64) (int64, error) {
	query := "DELETE FROM education WHERE id = $1"
	result, err := e.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete education: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("education with id %d not found", id)
	}

	return id, nil
}

// Get получает одну запись образования по ID
func (e *EducationRepo) Get(id int64) (models.Education, error) {
	query := "SELECT id, name, year, course, organization FROM education WHERE id = $1"
	
	var education models.Education
	err := e.db.QueryRow(query, id).Scan(
		&education.ID,
		&education.Name,
		&education.Year,
		&education.Course,
		&education.Organization,
	)
	
	if err == sql.ErrNoRows {
		return models.Education{}, fmt.Errorf("education with id %d not found", id)
	}
	if err != nil {
		return models.Education{}, fmt.Errorf("failed to get education: %w", err)
	}

	return education, nil
}

// List получает список записей образования с пагинацией, сортировкой и фильтрацией
func (e *EducationRepo) List(req entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Education], error) {
	baseQuery := "SELECT id, name, year, course, organization FROM education"

	queryParams := entityreqdecorator.BuildListQuery(
		req, baseQuery, e.isValidField,
	)

	// Получаем общее количество записей
	var total int
	err := e.db.QueryRow(queryParams.CountQuery, queryParams.CountParams...).Scan(&total)
	if err != nil {
		return entityreqdecorator.PagebleRs[models.Education]{}, fmt.Errorf("failed to count educations: %w", err)
	}

	// Получаем записи с учетом пагинации
	rows, err := e.db.Query(queryParams.SelectQuery, queryParams.SelectParams...)
	if err != nil {
		return entityreqdecorator.PagebleRs[models.Education]{}, fmt.Errorf("failed to query educations: %w", err)
	}
	defer rows.Close()

	var educations []models.Education
	for rows.Next() {
		var education models.Education
		err := rows.Scan(
			&education.ID,
			&education.Name,
			&education.Year,
			&education.Course,
			&education.Organization,
		)
		if err != nil {
			return entityreqdecorator.PagebleRs[models.Education]{}, fmt.Errorf("failed to scan education: %w", err)
		}
		educations = append(educations, education)
	}

	if err = rows.Err(); err != nil {
		return entityreqdecorator.PagebleRs[models.Education]{}, fmt.Errorf("rows error: %w", err)
	}

	return entityreqdecorator.PagebleRs[models.Education]{
		Total:   total,
		Content: educations,
		Page:    req.Page,
		Size:    req.Size,
		Sort:    req.Sort,
	}, nil
}

// Create создает новую запись образования
func (e *EducationRepo) Create(education models.Education) (models.Education, error) {
	query := `
		INSERT INTO education (name, year, course, organization)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, year, course, organization
	`

	var created models.Education
	err := e.db.QueryRow(
		query,
		education.Name,
		education.Year,
		education.Course,
		education.Organization,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Year,
		&created.Course,
		&created.Organization,
	)

	if err != nil {
		return models.Education{}, fmt.Errorf("failed to create education: %w", err)
	}

	return created, nil
}

// Update обновляет существующую запись образования
func (e *EducationRepo) Update(education models.Education) (models.Education, error) {
	query := `
		UPDATE education
		SET name = $2, year = $3, course = $4, organization = $5
		WHERE id = $1
		RETURNING id, name, year, course, organization
	`

	var updated models.Education
	err := e.db.QueryRow(
		query,
		education.ID,
		education.Name,
		education.Year,
		education.Course,
		education.Organization,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Year,
		&updated.Course,
		&updated.Organization,
	)

	if err == sql.ErrNoRows {
		return models.Education{}, fmt.Errorf("education with id %d not found", education.ID)
	}
	if err != nil {
		return models.Education{}, fmt.Errorf("failed to update education: %w", err)
	}

	return updated, nil
}

// isValidField проверяет, является ли поле валидным для сортировки и фильтрации
func (e *EducationRepo) isValidField(field string) bool {
	validFields := map[string]bool{
		"id":           true,
		"name":         true,
		"year":         true,
		"course":       true,
		"organization": true,
	}
	return validFields[field]
}
