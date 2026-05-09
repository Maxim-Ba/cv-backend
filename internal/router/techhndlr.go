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

// TechGet получает одну технологию по ID вместе с тегами
//
// @Summary      Получить технологию по ID
// @Tags         technologies
// @Produce      json
// @Param        techID  path  int  true  "ID технологии"
// @Success      200  {object}  dto.TechnologyWithTagsDTO
// @Failure      400  {object}  apierrors.APIError
// @Failure      404  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /tech/{techID} [get]
func (th *TechHandler) TechGet(w http.ResponseWriter, r *http.Request) {
	techIDStr := chi.URLParam(r, "techID")
	techID, err := strconv.ParseInt(techIDStr, 10, 64)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid technology ID")
		return
	}

	technology, err := th.service.GetWithTags(techID)
	if err != nil {
		if apierrors.IsNotFound(err) {
			apierrors.WriteError(w, http.StatusNotFound, "technology not found")
			return
		}
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(technology); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

// TechList получает список технологий с тегами
//
// @Summary      Список технологий
// @Tags         technologies
// @Produce      json
// @Param        page  query  int  false  "Номер страницы"
// @Param        size  query  int  false  "Размер страницы"
// @Success      200  {object}  dto.TechListResponse
// @Failure      500  {object}  apierrors.APIError
// @Router       /tech [get]
func (th *TechHandler) TechList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pagebleRq := entityreqdecorator.ParseQueryParams(queryParams)
	list, err := th.service.ListWithTags(pagebleRq)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

// TechCreate создает новую технологию
//
// @Summary      Создать технологию
// @Tags         technologies
// @Accept       json
// @Produce      json
// @Param        body  body  object{title=string,description=string,logoUrl=string}  true  "Данные технологии"
// @Success      201  {object}  dto.TechnologyDTO
// @Failure      400  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /tech [post]
func (th *TechHandler) TechCreate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		LogoUrl     string `json:"logoUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	technology := models.Technology{
		Title:       reqData.Title,
		Description: pgtype.Text{String: reqData.Description, Valid: reqData.Description != ""},
		LogoUrl:     pgtype.Text{String: reqData.LogoUrl, Valid: reqData.LogoUrl != ""},
	}

	created, err := th.service.Create(technology)
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

// TechDelete удаляет технологию
//
// @Summary      Удалить технологию(и)
// @Tags         technologies
// @Accept       json
// @Produce      json
// @Param        body  body  object{ids=[]int64}  true  "Список ID для удаления"
// @Success      200  {object}  dto.DeleteResponse
// @Failure      400  {object}  apierrors.APIError
// @Failure      404  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /tech [delete]
func (th *TechHandler) TechDelete(w http.ResponseWriter, r *http.Request) {
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
		if apierrors.IsNotFound(err) {
			apierrors.WriteError(w, http.StatusNotFound, "technology not found")
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

// TechUpdate обновляет технологию
//
// @Summary      Обновить технологию
// @Tags         technologies
// @Accept       json
// @Produce      json
// @Param        body  body  object{id=int64,title=string,description=string,logoUrl=string}  true  "Данные технологии"
// @Success      200  {object}  dto.TechnologyDTO
// @Failure      400  {object}  apierrors.APIError
// @Failure      404  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /tech [put]
func (th *TechHandler) TechUpdate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		ID          int64  `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		LogoUrl     string `json:"logoUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid request body")
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
		if apierrors.IsNotFound(err) {
			apierrors.WriteError(w, http.StatusNotFound, "technology not found")
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
