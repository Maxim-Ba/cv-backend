package services

import (
	"fmt"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

// EducationDeleter интерфейс для удаления записей образования
type EducationDeleter interface {
	DeleteList([]int64) ([]int64, error)
	Delete(id int64) (int64, error)
}

// EducationWriter интерфейс для создания и обновления записей образования
type EducationWriter interface {
	Create(models.Education) (models.Education, error)
	Update(models.Education) (models.Education, error)
}

// EducationReader интерфейс для чтения записей образования
type EducationReader interface {
	Get(id int64) (models.Education, error)
	List(entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Education], error)
}

// EducationManager объединяет все интерфейсы для работы с образованием
type EducationManager interface {
	EducationReader
	EducationWriter
	EducationDeleter
}

// EducationService сервис для работы с образованием
type EducationService struct {
	repo EducationManager
}

// NewEducationService создает новый экземпляр сервиса образования
func NewEducationService(repo EducationManager) *EducationService {
	return &EducationService{
		repo: repo,
	}
}

// DeleteList удаляет список записей образования по ID
func (s *EducationService) DeleteList(ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := s.repo.DeleteList(ids)
	if err != nil {
		return nil, fmt.Errorf("error deleting education list: %w", err)
	}
	return res, nil
}

// Delete удаляет одну запись образования по ID
func (s *EducationService) Delete(id int64) (int64, error) {
	if id == 0 {
		return 0, fmt.Errorf("invalid education ID: %d", id)
	}
	res, err := s.repo.Delete(id)
	if err != nil {
		return 0, fmt.Errorf("error deleting education: %w", err)
	}
	return res, nil
}

// Get получает одну запись образования по ID
func (s *EducationService) Get(id int64) (models.Education, error) {
	if id == 0 {
		return models.Education{}, fmt.Errorf("invalid education ID: %d", id)
	}
	res, err := s.repo.Get(id)
	if err != nil {
		return models.Education{}, fmt.Errorf("error getting education: %w", err)
	}
	return res, nil
}

// List получает список записей образования с пагинацией и фильтрацией
func (s *EducationService) List(r entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Education], error) {
	res, err := s.repo.List(r)
	if err != nil {
		return entityreqdecorator.PagebleRs[models.Education]{}, fmt.Errorf("error in getting list from Education repo: %w", err)
	}
	return res, nil
}

// Create создает новую запись образования
func (s *EducationService) Create(education models.Education) (models.Education, error) {
	if education.Course == "" || education.Organization == "" {
		return models.Education{}, fmt.Errorf("course and organization are required fields")
	}
	res, err := s.repo.Create(education)
	if err != nil {
		return models.Education{}, fmt.Errorf("error creating education: %w", err)
	}
	return res, nil
}

// Update обновляет существующую запись образования
func (s *EducationService) Update(education models.Education) (models.Education, error) {
	if education.ID == 0 {
		return models.Education{}, fmt.Errorf("invalid education ID: %d", education.ID)
	}
	if education.Course == "" || education.Organization == "" {
		return models.Education{}, fmt.Errorf("course and organization are required fields")
	}
	res, err := s.repo.Update(education)
	if err != nil {
		return models.Education{}, fmt.Errorf("error updating education: %w", err)
	}
	return res, nil
}
