package services

import (
	"fmt"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

type TagDeleter interface {
	DeleteList([]int64) ([]int64, error)
	Delete(id int64) (int64, error)
}
type TagWriter interface {
	Create(models.Tag) (models.Tag, error)
	Update(models.Tag) (models.Tag, error)
}

type TagReader interface {
	Get(id int64) (models.Tag, error)
	List(entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Tag], error)
}
type TagManager interface {
	TagReader
	TagWriter
	TagDeleter
}
type TagService struct {
	repo TagManager
}

func NewTagServise( repo TagManager) *TagService {
	return &TagService{
		repo: repo,
	}
}

func (s *TagService) DeleteList(ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := s.repo.DeleteList(ids)
	if err != nil {
		return nil, fmt.Errorf("error deleting tag list: %w", err)
	}
	return res, nil
}

func (s *TagService) Delete(id int64) (int64, error) {
	if id == 0 {
		return 0, fmt.Errorf("invalid tag ID: %d", id)
	}
	res, err := s.repo.Delete(id)
	if err != nil {
		return 0, fmt.Errorf("error deleting tag: %w", err)
	}
	return res, nil
}

// Get получает один тег по ID
func (s *TagService) Get(id int64) (models.Tag, error) {
	if id == 0 {
		return models.Tag{}, fmt.Errorf("invalid tag ID: %d", id)
	}
	res, err := s.repo.Get(id)
	if err != nil {
		return models.Tag{}, fmt.Errorf("error getting tag: %w", err)
	}
	return res, nil
}

// List получает список тегов с пагинацией и фильтрацией
func (s *TagService) List(r entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Tag], error) {
	res, err := s.repo.List(r)

	if err != nil {
		return entityreqdecorator.PagebleRs[models.Tag]{}, fmt.Errorf("error in getting list from tag repo: %w", err)
	}

	return res, nil
}

// Create создает новый тег
func (s *TagService) Create(tag models.Tag) (models.Tag, error) {
	if tag.Name == "" {
		return models.Tag{}, fmt.Errorf("tag name is required")
	}
	if tag.HexColor == "" {
		return models.Tag{}, fmt.Errorf("tag hex color is required")
	}
	res, err := s.repo.Create(tag)
	if err != nil {
		return models.Tag{}, fmt.Errorf("error creating tag: %w", err)
	}
	return res, nil
}

// Update обновляет существующий тег
func (s *TagService) Update(tag models.Tag) (models.Tag, error) {
	if tag.ID == 0 {
		return models.Tag{}, fmt.Errorf("invalid tag ID: %d", tag.ID)
	}
	if tag.Name == "" {
		return models.Tag{}, fmt.Errorf("tag name is required")
	}
	if tag.HexColor == "" {
		return models.Tag{}, fmt.Errorf("tag hex color is required")
	}
	res, err := s.repo.Update(tag)
	if err != nil {
		return models.Tag{}, fmt.Errorf("error updating tag: %w", err)
	}
	return res, nil
}
