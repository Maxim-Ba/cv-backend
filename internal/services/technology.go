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
	Create() (models.Technology, error)
	Update() (models.Technology, error)
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

	}
	return res, nil

}
func (s *TechService) Delete(id int64) (int64, error) {
	if id == 0 {
		return 0, nil
	}
	res, err := s.repo.Delete(id)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (s *TechService) List(r entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.Technology], error) {

	res, err := s.repo.List(r)

	if err != nil {
		return entityreqdecorator.PagebleRs[models.Technology]{}, fmt.Errorf("error in getting list from Tech repo %w ", err)
	}

	return res, nil
}
