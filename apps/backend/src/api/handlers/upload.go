package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/credfolio/apps/backend/src/services/profile"
	"github.com/google/uuid"
)

type UploadHandler struct {
	Service *profile.Service
}

func NewUploadHandler(s *profile.Service) *UploadHandler {
	return &UploadHandler{Service: s}
}

func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Parse multipart
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Dummy User ID for MVP
	userID := uuid.New()

	// Process
	if err := h.Service.ProcessUpload(r.Context(), userID, header.Filename, file); err != nil {
		slog.Error("upload processing failed", "error", err, "user_id", userID, "filename", header.Filename)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "processing",
		"user_id": userID.String(),
	})
}
