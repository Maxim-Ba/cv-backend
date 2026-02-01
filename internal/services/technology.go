package services

import (
	"fmt"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

type TechDeleter interface {
	DeleteList([]int64) ([]int64, error)
	Delete(id int64) (int64, error)
}
type TechWriter interface {
	Create(models.Technology) (models.Technology, error)
	Update(models.Technology) (models.Technology, error)
}

type TechReader interface {
	Get(id int64) (models.Technology, error)
	List(entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Technology], error)
}
type TechManager interface {
	TechReader
	TechWriter
	TechDeleter
}
type TechService struct {
	repo TechManager
}

func NewTechService(repo TechManager) *TechService {
	return &TechService{
		repo: repo,
	}
}
func (s *TechService) DeleteList(ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := s.repo.DeleteList(ids)
	if err != nil {
		return nil, fmt.Errorf("error deleting technology list: %w", err)
	}
	return res, nil
}

func (s *TechService) Delete(id int64) (int64, error) {
	if id == 0 {
		return 0, fmt.Errorf("invalid technology ID: %d", id)
	}
	res, err := s.repo.Delete(id)
	if err != nil {
		return 0, fmt.Errorf("error deleting technology: %w", err)
	}
	return res, nil
}

// Get получает одну технологию по ID
func (s *TechService) Get(id int64) (models.Technology, error) {
	if id == 0 {
		return models.Technology{}, fmt.Errorf("invalid technology ID: %d", id)
	}
	res, err := s.repo.Get(id)
	if err != nil {
		return models.Technology{}, fmt.Errorf("error getting technology: %w", err)
	}
	return res, nil
}

// List получает список технологий с пагинацией и фильтрацией
func (s *TechService) List(r entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Technology], error) {
	res, err := s.repo.List(r)

	if err != nil {
		return entityreqdecorator.PagebleRs[models.Technology]{}, fmt.Errorf("error in getting list from Tech repo: %w", err)
	}

	return res, nil
}

// Create создает новую технологию
func (s *TechService) Create(technology models.Technology) (models.Technology, error) {
	if technology.Title == "" {
		return models.Technology{}, fmt.Errorf("technology title is required")
	}
	res, err := s.repo.Create(technology)
	if err != nil {
		return models.Technology{}, fmt.Errorf("error creating technology: %w", err)
	}
	return res, nil
}

// Update обновляет существующую технологию
func (s *TechService) Update(technology models.Technology) (models.Technology, error) {
	if technology.ID == 0 {
		return models.Technology{}, fmt.Errorf("invalid technology ID: %d", technology.ID)
	}
	if technology.Title == "" {
		return models.Technology{}, fmt.Errorf("technology title is required")
	}
	res, err := s.repo.Update(technology)
	if err != nil {
		return models.Technology{}, fmt.Errorf("error updating technology: %w", err)
	}
	return res, nil
}
