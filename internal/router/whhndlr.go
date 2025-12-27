package router

import (
	"net/http"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

type WorkHistoryDeleter interface {
	DeleteList([]int64) ([]int64, error)
	Delete() (int64, error)
}
type WorkHistoryWriter interface {
	Create() (models.WorkHistory, error)
	Update() (models.WorkHistory, error)
}

type WorkHistoryReader interface {
	Get() (models.WorkHistory, error)
	List() (entityreqdecorator.PagebleRs[models.WorkHistory], error)
}

func WorkHistoryGet(w http.ResponseWriter, r *http.Request)  {}
func WorkHistoryList(w http.ResponseWriter, r *http.Request) {}

func WorkHistoryCreate(w http.ResponseWriter, r *http.Request) {}

func WorkHistoryDelete(w http.ResponseWriter, r *http.Request) {}

func WorkHistoryUpdate(w http.ResponseWriter, r *http.Request) {}
