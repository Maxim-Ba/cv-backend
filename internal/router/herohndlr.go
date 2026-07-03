package router

import (
	"encoding/json"
	"net/http"

	_ "github.com/Maxim-Ba/cv-backend/internal/models/dto"
	"github.com/Maxim-Ba/cv-backend/internal/services"
	"github.com/Maxim-Ba/cv-backend/pkg/apierrors"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
)

// HeroHandler хендлер hero-секции
type HeroHandler struct {
	service *services.ProfileService
}

// NewHeroHandler создает HeroHandler
func NewHeroHandler(ps *services.ProfileService) *HeroHandler {
	return &HeroHandler{service: ps}
}

// HeroGet возвращает данные hero-секции
//
// @Summary      Получить hero-секцию
// @Tags         hero
// @Produce      json
// @Success      200  {object}  dto.HeroDTO
// @Failure      404  {object}  apierrors.APIError
// @Failure      500  {object}  apierrors.APIError
// @Router       /hero [get]
func (h *HeroHandler) HeroGet(w http.ResponseWriter, r *http.Request) {
	locale := i18n.FromContext(r.Context())
	hero, err := h.service.GetHero(locale)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hero); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "failed to encode response")
	}
}
