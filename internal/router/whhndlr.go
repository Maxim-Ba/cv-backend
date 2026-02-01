package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	"github.com/Maxim-Ba/cv-backend/internal/services"
	entityreqdecorator "github.com/Maxim-Ba/cv-backend/pkg/entity-req-decorator"
)

// WorkHistoryHandler хендлер для работы с историей работы
type WorkHistoryHandler struct {
	service *services.WorkHistoryService
}

// NewWorkHistoryHandler создает новый экземпляр хендлера истории работы
func NewWorkHistoryHandler(whs *services.WorkHistoryService) *WorkHistoryHandler {
	return &WorkHistoryHandler{
		service: whs,
	}
}

// WorkHistoryGet получает одну запись истории работы по ID
func (wh *WorkHistoryHandler) WorkHistoryGet(w http.ResponseWriter, r *http.Request) {
	whIDStr := chi.URLParam(r, "whID")
	whID, err := strconv.ParseInt(whIDStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid work history ID",
		})
		return
	}

	workHistory, err := wh.service.Get(whID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(workHistory); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// WorkHistoryList получает список записей истории работы
func (wh *WorkHistoryHandler) WorkHistoryList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	list, err := wh.service.List(pagebleRq)

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

// WorkHistoryCreate создает новую запись истории работы
func (wh *WorkHistoryHandler) WorkHistoryCreate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Name        string   `json:"name"`
		About       string   `json:"about"`
		LogoUrl     []byte   `json:"logoUrl"`
		PeriodStart string   `json:"periodStart"`
		PeriodEnd   string   `json:"periodEnd"`
		WhatIDid    []string `json:"whatIDid"`
		Projects    []string `json:"projects"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Парсинг дат
	var periodStart, periodEnd pgtype.Date
	if reqData.PeriodStart != "" {
		t, err := time.Parse("2006-01-02", reqData.PeriodStart)
		if err == nil {
			periodStart = pgtype.Date{Time: t, Valid: true}
		}
	}
	if reqData.PeriodEnd != "" {
		t, err := time.Parse("2006-01-02", reqData.PeriodEnd)
		if err == nil {
			periodEnd = pgtype.Date{Time: t, Valid: true}
		}
	}

	workHistory := models.WorkHistory{
		Name:        reqData.Name,
		About:       reqData.About,
		LogoUrl:     reqData.LogoUrl,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		WhatIDid:    reqData.WhatIDid,
		Projects:    reqData.Projects,
	}

	created, err := wh.service.Create(workHistory)
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

// WorkHistoryDelete удаляет запись истории работы
func (wh *WorkHistoryHandler) WorkHistoryDelete(w http.ResponseWriter, r *http.Request) {
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
		deletedID, delErr := wh.service.Delete(deleteReq.IDs[0])
		if delErr != nil {
			err = delErr
		} else {
			deletedIDs = []int64{deletedID}
		}
	} else if len(deleteReq.IDs) > 1 {
		deletedIDs, err = wh.service.DeleteList(deleteReq.IDs)
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

// WorkHistoryUpdate обновляет запись истории работы
func (wh *WorkHistoryHandler) WorkHistoryUpdate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		ID          int64    `json:"id"`
		Name        string   `json:"name"`
		About       string   `json:"about"`
		LogoUrl     []byte   `json:"logoUrl"`
		PeriodStart string   `json:"periodStart"`
		PeriodEnd   string   `json:"periodEnd"`
		WhatIDid    []string `json:"whatIDid"`
		Projects    []string `json:"projects"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Парсинг дат
	var periodStart, periodEnd pgtype.Date
	if reqData.PeriodStart != "" {
		t, err := time.Parse("2006-01-02", reqData.PeriodStart)
		if err == nil {
			periodStart = pgtype.Date{Time: t, Valid: true}
		}
	}
	if reqData.PeriodEnd != "" {
		t, err := time.Parse("2006-01-02", reqData.PeriodEnd)
		if err == nil {
			periodEnd = pgtype.Date{Time: t, Valid: true}
		}
	}

	workHistory := models.WorkHistory{
		ID:          reqData.ID,
		Name:        reqData.Name,
		About:       reqData.About,
		LogoUrl:     reqData.LogoUrl,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		WhatIDid:    reqData.WhatIDid,
		Projects:    reqData.Projects,
	}

	updated, err := wh.service.Update(workHistory)
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
