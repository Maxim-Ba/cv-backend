package services

import (
	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

type TagDeleter interface {
	DeleteList([]int64) ([]int64, error)
	Delete(id int64) (int64, error)
}
type TagWriter interface {
	Create() (models.Tag, error)
	Update() (models.Tag, error)
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

	}
	return res, nil

}
func (s *TagService) Delete(id int64) (int64, error) {
	if id == 0 {
		return 0, nil
	}
	res, err := s.repo.Delete(id)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *TagService) List(r entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Tag], error) {
	if r.Page == 0 {
		return entityreqdecorator.PagebleRs[models.Tag]{}, nil
	}
	if r.Size == 0 {
		return entityreqdecorator.PagebleRs[models.Tag]{}, nil
	}

	res, err := s.repo.List(r)

	if err != nil {
		return entityreqdecorator.PagebleRs[models.Tag]{}, err
	}

	return res, nil
}
