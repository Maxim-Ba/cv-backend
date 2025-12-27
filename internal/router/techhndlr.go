package router

import (
	"net/http"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

type TechnologyDeleter interface {
	DeleteList([]int64) ([]int64, error)
	Delete() (int64, error)
}
type TechnologyWriter interface {
	Create() (models.Technology, error)
	Update() (models.Technology, error)
}

type TechnologyReader interface {
	Get() (models.Technology, error)
	List() (entityreqdecorator.PagebleRs[models.Technology], error)
}

func TechGet(w http.ResponseWriter, r *http.Request)  {}
func TechList(w http.ResponseWriter, r *http.Request) {}

func TechCreate(w http.ResponseWriter, r *http.Request) {}

func TechDelete(w http.ResponseWriter, r *http.Request) {}

func TechUpdate(w http.ResponseWriter, r *http.Request) {}
