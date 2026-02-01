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

// EducationHandler хендлер для работы с образованием
type EducationHandler struct {
	service *services.EducationService
}

// NewEducationHandler создает новый экземпляр хендлера образования
func NewEducationHandler(es *services.EducationService) *EducationHandler {
	return &EducationHandler{
		service: es,
	}
}

// EducationGet получает одну запись образования по ID
func (eh *EducationHandler) EducationGet(w http.ResponseWriter, r *http.Request) {
	eduIDStr := chi.URLParam(r, "eduID")
	eduID, err := strconv.ParseInt(eduIDStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid education ID",
		})
		return
	}

	education, err := eh.service.Get(eduID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(education); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// EducationList получает список записей образования
func (eh *EducationHandler) EducationList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	list, err := eh.service.List(pagebleRq)

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

// EducationCreate создает новую запись образования
func (eh *EducationHandler) EducationCreate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Name         string `json:"name"`
		Year         int32  `json:"year"`
		Course       string `json:"course"`
		Organization string `json:"organization"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	education := models.Education{
		Name:         pgtype.Text{String: reqData.Name, Valid: reqData.Name != ""},
		Year:         reqData.Year,
		Course:       reqData.Course,
		Organization: reqData.Organization,
	}

	created, err := eh.service.Create(education)
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

// EducationDelete удаляет запись образования
func (eh *EducationHandler) EducationDelete(w http.ResponseWriter, r *http.Request) {
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
		deletedID, delErr := eh.service.Delete(deleteReq.IDs[0])
		if delErr != nil {
			err = delErr
		} else {
			deletedIDs = []int64{deletedID}
		}
	} else if len(deleteReq.IDs) > 1 {
		deletedIDs, err = eh.service.DeleteList(deleteReq.IDs)
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

// EducationUpdate обновляет запись образования
func (eh *EducationHandler) EducationUpdate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		Year         int32  `json:"year"`
		Course       string `json:"course"`
		Organization string `json:"organization"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	education := models.Education{
		ID:           reqData.ID,
		Name:         pgtype.Text{String: reqData.Name, Valid: reqData.Name != ""},
		Year:         reqData.Year,
		Course:       reqData.Course,
		Organization: reqData.Organization,
	}

	updated, err := eh.service.Update(education)
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
