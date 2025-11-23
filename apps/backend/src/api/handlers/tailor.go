package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/credfolio/apps/backend/src/services/profile"
	"github.com/google/uuid"
)

type TailorHandler struct {
	Service *profile.Service
}

func NewTailorHandler(s *profile.Service) *TailorHandler {
	return &TailorHandler{Service: s}
}

type TailorRequest struct {
	JobDescription string `json:"job_description"`
	UserID         string `json:"user_id"`
}

func (h *TailorHandler) TailorProfile(w http.ResponseWriter, r *http.Request) {
	var req TailorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.JobDescription == "" || req.UserID == "" {
		http.Error(w, "job_description and user_id required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	res, err := h.Service.TailorProfile(r.Context(), userID, req.JobDescription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

