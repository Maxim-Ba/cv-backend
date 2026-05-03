package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	_ "github.com/Maxim-Ba/cv-backend/internal/models/dto"
	models "github.com/Maxim-Ba/cv-backend/internal/models/gen"
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

// TagGet получает один тег по ID
//
// @Summary      Получить тег по ID
// @Tags         tags
// @Produce      json
// @Param        tagID  path  int  true  "ID тега"
// @Success      200  {object}  dto.TagDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tag/{tagID} [get]
func (th *TagHandler) TagGet(w http.ResponseWriter, r *http.Request) {
	tagIDStr := chi.URLParam(r, "tagID")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid tag ID",
		})
		return
	}

	tag, err := th.service.Get(tagID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tag); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
// TagList получает список тегов с пагинацией
//
// @Summary      Список тегов
// @Tags         tags
// @Produce      json
// @Param        page  query  int  false  "Номер страницы"
// @Param        size  query  int  false  "Размер страницы"
// @Success      200  {object}  dto.TagListResponse
// @Failure      500  {object}  map[string]string
// @Router       /tag [get]
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

// TagCreate создает новый тег
//
// @Summary      Создать тег
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        body  body  object{name=string,hexColor=string}  true  "Данные тега"
// @Success      201  {object}  dto.TagDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tag [post]
func (th *TagHandler) TagCreate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Name     string `json:"name"`
		HexColor string `json:"hexColor"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	tag := models.Tag{
		Name:     reqData.Name,
		HexColor: reqData.HexColor,
	}

	created, err := th.service.Create(tag)
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

// TagDelete удаляет тег
//
// @Summary      Удалить тег(и)
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        body  body  object{ids=[]int64}  true  "Список ID для удаления"
// @Success      200  {object}  dto.DeleteResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tag [delete]
func (th *TagHandler) TagDelete(w http.ResponseWriter, r *http.Request) {
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

// TagUpdate обновляет тег
//
// @Summary      Обновить тег
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        body  body  object{id=int64,name=string,hexColor=string}  true  "Данные тега"
// @Success      200  {object}  dto.TagDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tag [put]
func (th *TagHandler) TagUpdate(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		HexColor string `json:"hexColor"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	tag := models.Tag{
		ID:       reqData.ID,
		Name:     reqData.Name,
		HexColor: reqData.HexColor,
	}

	updated, err := th.service.Update(tag)
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
