package router

import (
	"encoding/json"
	"net/http"

	"github.com/Maxim-Ba/cv-backend/internal/services"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

type TagHandler struct {
	service services.TagService
}

func NewTagHandler(ts services.TagService) *TagHandler {
	return &TagHandler{
		service: ts,
	}
}

func (th *TagHandler) TagGet(w http.ResponseWriter, r *http.Request) {

}
func (th *TagHandler) TagList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	list, err := th.service.List(pagebleRq)

	if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": err.Error(),
        })
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(list); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

func (th *TagHandler) TagCreate(w http.ResponseWriter, r *http.Request) {}

func (th *TagHandler) TagDelete(w http.ResponseWriter, r *http.Request) {}

func (th *TagHandler) TagUpdate(w http.ResponseWriter, r *http.Request) {}
