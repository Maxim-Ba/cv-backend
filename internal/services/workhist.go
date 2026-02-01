package services

import (
	"fmt"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

// WorkHistoryDeleter интерфейс для удаления записей истории работы
type WorkHistoryDeleter interface {
	DeleteList([]int64) ([]int64, error)
	Delete(id int64) (int64, error)
}

// WorkHistoryWriter интерфейс для создания и обновления записей истории работы
type WorkHistoryWriter interface {
	Create(models.WorkHistory) (models.WorkHistory, error)
	Update(models.WorkHistory) (models.WorkHistory, error)
}

// WorkHistoryReader интерфейс для чтения записей истории работы
type WorkHistoryReader interface {
	Get(id int64) (models.WorkHistory, error)
	List(entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.WorkHistory], error)
}

// WorkHistoryManager объединяет все интерфейсы для работы с историей работы
type WorkHistoryManager interface {
	WorkHistoryReader
	WorkHistoryWriter
	WorkHistoryDeleter
}

// WorkHistoryService сервис для работы с историей работы
type WorkHistoryService struct {
	repo WorkHistoryManager
}

// NewWorkHistoryService создает новый экземпляр сервиса истории работы
func NewWorkHistoryService(repo WorkHistoryManager) *WorkHistoryService {
	return &WorkHistoryService{
		repo: repo,
	}
}

// DeleteList удаляет список записей истории работы по ID
func (s *WorkHistoryService) DeleteList(ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := s.repo.DeleteList(ids)
	if err != nil {
		return nil, fmt.Errorf("error deleting work history list: %w", err)
	}
	return res, nil
}

// Delete удаляет одну запись истории работы по ID
func (s *WorkHistoryService) Delete(id int64) (int64, error) {
	if id == 0 {
		return 0, fmt.Errorf("invalid work history ID: %d", id)
	}
	res, err := s.repo.Delete(id)
	if err != nil {
		return 0, fmt.Errorf("error deleting work history: %w", err)
	}
	return res, nil
}

// Get получает одну запись истории работы по ID
func (s *WorkHistoryService) Get(id int64) (models.WorkHistory, error) {
	if id == 0 {
		return models.WorkHistory{}, fmt.Errorf("invalid work history ID: %d", id)
	}
	res, err := s.repo.Get(id)
	if err != nil {
		return models.WorkHistory{}, fmt.Errorf("error getting work history: %w", err)
	}
	return res, nil
}

// List получает список записей истории работы с пагинацией и фильтрацией
func (s *WorkHistoryService) List(r entityreqdecorator.PagebleRq) (entityreqdecorator.PagebleRs[models.WorkHistory], error) {
	res, err := s.repo.List(r)
	if err != nil {
		return entityreqdecorator.PagebleRs[models.WorkHistory]{}, fmt.Errorf("error in getting list from WorkHistory repo: %w", err)
	}
	return res, nil
}

// Create создает новую запись истории работы
func (s *WorkHistoryService) Create(workHistory models.WorkHistory) (models.WorkHistory, error) {
	if workHistory.Name == "" || workHistory.About == "" {
		return models.WorkHistory{}, fmt.Errorf("name and about are required fields")
	}
	res, err := s.repo.Create(workHistory)
	if err != nil {
		return models.WorkHistory{}, fmt.Errorf("error creating work history: %w", err)
	}
	return res, nil
}

// Update обновляет существующую запись истории работы
func (s *WorkHistoryService) Update(workHistory models.WorkHistory) (models.WorkHistory, error) {
	if workHistory.ID == 0 {
		return models.WorkHistory{}, fmt.Errorf("invalid work history ID: %d", workHistory.ID)
	}
	if workHistory.Name == "" || workHistory.About == "" {
		return models.WorkHistory{}, fmt.Errorf("name and about are required fields")
	}
	res, err := s.repo.Update(workHistory)
	if err != nil {
		return models.WorkHistory{}, fmt.Errorf("error updating work history: %w", err)
	}
	return res, nil
}
