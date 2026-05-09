package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	_ "github.com/Maxim-Ba/cv-backend/internal/models/dto"
	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
	"github.com/Maxim-Ba/cv-backend/internal/services"
	"github.com/Maxim-Ba/cv-backend/pkg/apierrors"
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
//
// @Summary      Получить запись образования по ID
// @Tags         education
// @Produce      json
// @Param        eduID  path  int  true  "ID записи образования"
// @Success      200  {object}  dto.EducationDTO
// @Failure      400  {object}  apierrors.APIError
// @Failure      404  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /edu/{eduID} [get]
func (eh *EducationHandler) EducationGet(w http.ResponseWriter, r *http.Request) {
	eduIDStr := chi.URLParam(r, "eduID")
	eduID, err := strconv.ParseInt(eduIDStr, 10, 64)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid education ID")
		return
	}

	education, err := eh.service.Get(eduID)
	if err != nil {
		if apierrors.IsNotFound(err) {
			apierrors.WriteError(w, http.StatusNotFound, "education not found")
			return
		}
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(education); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

// EducationList получает список записей образования
//
// @Summary      Список записей образования
// @Tags         education
// @Produce      json
// @Param        page  query  int  false  "Номер страницы"
// @Param        size  query  int  false  "Размер страницы"
// @Success      200  {object}  dto.EducationListResponse
// @Failure      500  {object}  apierrors.APIError
// @Router       /edu [get]
func (eh *EducationHandler) EducationList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	list, err := eh.service.List(pagebleRq)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

// EducationCreate создает новую запись образования
//
// @Summary      Создать запись образования
// @Tags         education
// @Accept       json
// @Produce      json
// @Param        body  body  object{name=string,year=int32,course=string,organization=string}  true  "Данные записи"
// @Success      201  {object}  dto.EducationDTO
// @Failure      400  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /edu [post]
func (eh *EducationHandler) EducationCreate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Name         string `json:"name"`
		Year         int32  `json:"year"`
		Course       string `json:"course"`
		Organization string `json:"organization"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid request body")
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
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

// EducationDelete удаляет запись образования
//
// @Summary      Удалить запись(и) образования
// @Tags         education
// @Accept       json
// @Produce      json
// @Param        body  body  object{ids=[]int64}  true  "Список ID для удаления"
// @Success      200  {object}  dto.DeleteResponse
// @Failure      400  {object}  apierrors.APIError
// @Failure      404  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /edu [delete]
func (eh *EducationHandler) EducationDelete(w http.ResponseWriter, r *http.Request) {
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
		if apierrors.IsNotFound(err) {
			apierrors.WriteError(w, http.StatusNotFound, "education not found")
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

// EducationUpdate обновляет запись образования
//
// @Summary      Обновить запись образования
// @Tags         education
// @Accept       json
// @Produce      json
// @Param        body  body  object{id=int64,name=string,year=int32,course=string,organization=string}  true  "Данные записи"
// @Success      200  {object}  dto.EducationDTO
// @Failure      400  {object}  apierrors.APIError
// @Failure      404  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /edu [put]
func (eh *EducationHandler) EducationUpdate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		Year         int32  `json:"year"`
		Course       string `json:"course"`
		Organization string `json:"organization"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid request body")
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
		if apierrors.IsNotFound(err) {
			apierrors.WriteError(w, http.StatusNotFound, "education not found")
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
