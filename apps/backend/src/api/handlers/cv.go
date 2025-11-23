package handlers

import (
	"net/http"

	"github.com/credfolio/apps/backend/src/services/generator"
	"github.com/credfolio/apps/backend/src/services/profile"
	"github.com/google/uuid"
)

type CVHandler struct {
	Service   *profile.Service
	Generator *generator.CVGenerator
}

func NewCVHandler(s *profile.Service, g *generator.CVGenerator) *CVHandler {
	return &CVHandler{Service: s, Generator: g}
}

func (h *CVHandler) DownloadCV(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	prof, err := h.Service.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	pdfBytes, err := h.Generator.Generate(prof, nil)
	if err != nil {
		http.Error(w, "generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="credfolio_cv.pdf"`)
	_, _ = w.Write(pdfBytes)
}
