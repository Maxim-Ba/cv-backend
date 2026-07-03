package router

import (
	"encoding/json"
	"net/http"

	_ "github.com/Maxim-Ba/cv-backend/internal/models/dto"
	"github.com/Maxim-Ba/cv-backend/internal/services"
	"github.com/Maxim-Ba/cv-backend/pkg/apierrors"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
)

// AboutMeHandler хендлер для секции «О себе»
type AboutMeHandler struct {
	service *services.ProfileService
}

// NewAboutMeHandler создает новый экземпляр хендлера секции «О себе»
func NewAboutMeHandler(ps *services.ProfileService) *AboutMeHandler {
	return &AboutMeHandler{service: ps}
}

// AboutMeGet возвращает данные секции «О себе»
func (h *AboutMeHandler) AboutMeGet(w http.ResponseWriter, r *http.Request) {
	locale := i18n.FromContext(r.Context())
	aboutMe, err := h.service.GetAboutMe(locale)
	if err != nil {
		if apierrors.IsNotFound(err) {
			apierrors.WriteError(w, http.StatusNotFound, "about me not found")
			return
		}
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(aboutMe); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

// AboutMeUpdate обновляет данные секции «О себе»
func (h *AboutMeHandler) AboutMeUpdate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BioParagraphs []string `json:"bioParagraphs"`
		Note          string   `json:"note"`
		Hobbies       string   `json:"hobbies"`
		TechnologyIds []int64  `json:"technologyIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := services.UpdateAboutMeInput{
		About:         i18n.FromLegacyText(services.JoinBioParagraphs(req.BioParagraphs)),
		Note:          i18n.FromLegacyText(req.Note),
		Hobbies:       i18n.FromLegacyText(req.Hobbies),
		TechnologyIDs: req.TechnologyIds,
	}
	if err := h.service.UpdateAboutMe(input); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	locale := i18n.FromContext(r.Context())
	aboutMe, err := h.service.GetAboutMe(locale)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(aboutMe); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}
