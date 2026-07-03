package router

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Maxim-Ba/cv-backend/internal/services"
	"github.com/Maxim-Ba/cv-backend/pkg/apierrors"
	"github.com/Maxim-Ba/cv-backend/pkg/i18n"
)

// PDFHandler хендлер для генерации и скачивания CV в PDF
type PDFHandler struct {
	svc *services.PDFService
}

// newPDFHandler создает новый экземпляр PDFHandler
func newPDFHandler(svc *services.PDFService) *PDFHandler {
	return &PDFHandler{svc: svc}
}

// DownloadCV godoc
// @Summary      Скачать CV в PDF
// @Description  Генерирует актуальное резюме из данных БД и возвращает PDF-файл
// @Tags         cv
// @Produce      application/pdf
// @Success      200  {file}    binary
// @Failure      500  {object}  apierrors.APIError
// @Router       /download-cv [get]
func (h *PDFHandler) DownloadCV(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.GenerateCV(i18n.ResolveLocale(r))
	if err != nil {
		slog.Error("failed to generate PDF", "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, "Не удалось сгенерировать PDF")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="CV_Balashov_Maxim.pdf"`)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data) //nolint:errcheck
}
