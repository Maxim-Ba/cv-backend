package router

import (
	"net/http"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

type EducationDeleter interface {
	DeleteList([]int64) ([]int64, error)
	Delete() (int64, error)
}
type EducationWriter interface {
	Create() (models.Education, error)
	Update() (models.Education, error)
}

type EducationReader interface {
	Get() (models.Education, error)
	List() (entityreqdecorator.PagebleRs[models.Education], error)
}
func EducationGet(w http.ResponseWriter, r *http.Request)  {}
func EducationList(w http.ResponseWriter, r *http.Request) {}

func EducationCreate(w http.ResponseWriter, r *http.Request) {}

func EducationDelete(w http.ResponseWriter, r *http.Request) {}

func EducationUpdate(w http.ResponseWriter, r *http.Request) {}
