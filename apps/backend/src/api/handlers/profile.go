package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/credfolio/apps/backend/src/services/profile"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	Service *profile.Service
}

func NewProfileHandler(s *profile.Service) *ProfileHandler {
	return &ProfileHandler{Service: s}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Dummy User ID for MVP
	// In real app, get from context
	// For now, hardcode or create new if not found?
	// To make it work with upload, we need the SAME ID.
	// Since we don't have auth/session, we can't persist ID across requests easily without client storage.
	// For MVP demo, we might want to pass ID or just use a fixed ID for dev.

	// Let's try to get ID from query param or header, else hardcoded.
	// Or we rely on single-user mode.
	userIDStr := r.URL.Query().Get("user_id")
	var userID uuid.UUID
	var err error

	if userIDStr != "" {
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
	} else {
		// Fallback for dev: hardcoded specific UUID or similar if we want persistence?
		// But ID is generated randomly in Upload handler.
		// We should probably return the ID in Upload response and Client sends it back.
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	profile, err := h.Service.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(profile)
}
