package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	_ "github.com/Maxim-Ba/cv-backend/internal/models/dto"
	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	"github.com/Maxim-Ba/cv-backend/internal/services"
	"github.com/Maxim-Ba/cv-backend/pkg/apierrors"
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

// WorkHistoryGet получает одну запись истории работы по ID с технологиями
//
// @Summary      Получить запись истории работы по ID
// @Tags         work-history
// @Produce      json
// @Param        whID  path  int  true  "ID записи истории работы"
// @Success      200  {object}  dto.WorkHistoryWithTechnologiesDTO
// @Failure      400  {object}  apierrors.APIError
// @Failure      404  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /wh/{whID} [get]
func (wh *WorkHistoryHandler) WorkHistoryGet(w http.ResponseWriter, r *http.Request) {
	whIDStr := chi.URLParam(r, "whID")
	whID, err := strconv.ParseInt(whIDStr, 10, 64)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid work history ID")
		return
	}

	workHistory, err := wh.service.GetWithTechnologies(whID)
	if err != nil {
		if apierrors.IsNotFound(err) {
			apierrors.WriteError(w, http.StatusNotFound, "work history not found")
			return
		}
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(workHistory); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

// WorkHistoryList получает список записей истории работы с технологиями
//
// @Summary      Список истории работы
// @Tags         work-history
// @Produce      json
// @Param        page  query  int  false  "Номер страницы"
// @Param        size  query  int  false  "Размер страницы"
// @Success      200  {object}  dto.WorkHistoryListResponse
// @Failure      500  {object}  apierrors.APIError
// @Router       /wh [get]
func (wh *WorkHistoryHandler) WorkHistoryList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	list, err := wh.service.ListWithTechnologies(pagebleRq)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

// WorkHistoryCreate создает новую запись истории работы
//
// @Summary      Создать запись истории работы
// @Tags         work-history
// @Accept       json
// @Produce      json
// @Param        body  body  object{name=string,about=string,logoUrl=string,periodStart=string,periodEnd=string,whatIDid=[]string,projects=[]string}  true  "Данные записи"
// @Success      201  {object}  dto.WorkHistoryDTO
// @Failure      400  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /wh [post]
func (wh *WorkHistoryHandler) WorkHistoryCreate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Name        string   `json:"name"`
		About       string   `json:"about"`
		LogoUrl     string   `json:"logoUrl"`
		PeriodStart string   `json:"periodStart"`
		PeriodEnd   string   `json:"periodEnd"`
		WhatIDid    []string `json:"whatIDid"`
		Projects    []string `json:"projects"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid request body")
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
		LogoUrl:     pgtype.Text{String: reqData.LogoUrl, Valid: reqData.LogoUrl != ""},
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		WhatIDid:    reqData.WhatIDid,
		Projects:    reqData.Projects,
	}

	created, err := wh.service.Create(workHistory)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

// WorkHistoryDelete удаляет запись истории работы
//
// @Summary      Удалить запись(и) истории работы
// @Tags         work-history
// @Accept       json
// @Produce      json
// @Param        body  body  object{ids=[]int64}  true  "Список ID для удаления"
// @Success      200  {object}  dto.DeleteResponse
// @Failure      400  {object}  apierrors.APIError
// @Failure      404  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /wh [delete]
func (wh *WorkHistoryHandler) WorkHistoryDelete(w http.ResponseWriter, r *http.Request) {
	var deleteReq struct {
		IDs []int64 `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deleteReq); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid request body")
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
		if apierrors.IsNotFound(err) {
			apierrors.WriteError(w, http.StatusNotFound, "work history not found")
			return
		}
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
		"deleted_ids": deletedIDs,
		"count":       len(deletedIDs),
	})
}

// WorkHistoryUpdate обновляет запись истории работы
//
// @Summary      Обновить запись истории работы
// @Tags         work-history
// @Accept       json
// @Produce      json
// @Param        body  body  object{id=int64,name=string,about=string,logoUrl=string,periodStart=string,periodEnd=string,whatIDid=[]string,projects=[]string}  true  "Данные записи"
// @Success      200  {object}  dto.WorkHistoryDTO
// @Failure      400  {object}  apierrors.APIError
// @Failure      404  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /wh [put]
func (wh *WorkHistoryHandler) WorkHistoryUpdate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		ID          int64    `json:"id"`
		Name        string   `json:"name"`
		About       string   `json:"about"`
		LogoUrl     string   `json:"logoUrl"`
		PeriodStart string   `json:"periodStart"`
		PeriodEnd   string   `json:"periodEnd"`
		WhatIDid    []string `json:"whatIDid"`
		Projects    []string `json:"projects"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid request body")
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
		LogoUrl:     pgtype.Text{String: reqData.LogoUrl, Valid: reqData.LogoUrl != ""},
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		WhatIDid:    reqData.WhatIDid,
		Projects:    reqData.Projects,
	}

	updated, err := wh.service.Update(workHistory)
	if err != nil {
		if apierrors.IsNotFound(err) {
			apierrors.WriteError(w, http.StatusNotFound, "work history not found")
			return
		}
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}
