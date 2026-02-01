package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	"github.com/Maxim-Ba/cv-backend/internal/services"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

// TechHandler хендлер для работы с технологиями
type TechHandler struct {
	service *services.TechService
}

// NewTechHandler создает новый экземпляр хендлера технологий
func NewTechHandler(ts *services.TechService) *TechHandler {
	return &TechHandler{
		service: ts,
	}
}

// TechGet получает одну технологию по ID
func (th *TechHandler) TechGet(w http.ResponseWriter, r *http.Request) {
	techIDStr := chi.URLParam(r, "techID")
	techID, err := strconv.ParseInt(techIDStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid technology ID",
		})
		return
	}

	technology, err := th.service.Get(techID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(technology); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// TechList получает список технологий
func (th *TechHandler) TechList(w http.ResponseWriter, r *http.Request) {
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

// TechCreate создает новую технологию
func (th *TechHandler) TechCreate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		LogoUrl     string `json:"logoUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	technology := models.Technology{
		Title:       reqData.Title,
		Description: pgtype.Text{String: reqData.Description, Valid: reqData.Description != ""},
		LogoUrl:     pgtype.Text{String: reqData.LogoUrl, Valid: reqData.LogoUrl != ""},
	}

	created, err := th.service.Create(technology)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// TechDelete удаляет технологию
func (th *TechHandler) TechDelete(w http.ResponseWriter, r *http.Request) {
	var deleteReq struct {
		IDs []int64 `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deleteReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	var deletedIDs []int64
	var err error

	if len(deleteReq.IDs) == 1 {
		deletedID, delErr := th.service.Delete(deleteReq.IDs[0])
		if delErr != nil {
			err = delErr
		} else {
			deletedIDs = []int64{deletedID}
		}
	} else if len(deleteReq.IDs) > 1 {
		deletedIDs, err = th.service.DeleteList(deleteReq.IDs)
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deleted_ids": deletedIDs,
		"count":       len(deletedIDs),
	})
}

// TechUpdate обновляет технологию
func (th *TechHandler) TechUpdate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		ID          int64  `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		LogoUrl     string `json:"logoUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	technology := models.Technology{
		ID:          reqData.ID,
		Title:       reqData.Title,
		Description: pgtype.Text{String: reqData.Description, Valid: reqData.Description != ""},
		LogoUrl:     pgtype.Text{String: reqData.LogoUrl, Valid: reqData.LogoUrl != ""},
	}

	updated, err := th.service.Update(technology)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
