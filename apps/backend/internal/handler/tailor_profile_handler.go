package handler

import (
	"encoding/json"
	"net/http"

	"github.com/credfolio/apps/backend/api/generated"
	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/google/uuid"
)

// TailorProfile implements generated.ServerInterface
func (a *API) TailorProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var body generated.TailorProfileJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate job description
	if body.JobDescription == nil || *body.JobDescription == "" {
		http.Error(w, "job description is required", http.StatusBadRequest)
		return
	}

	// Call tailoring service
	jobMatch, err := a.TailoringService.TailorProfileToJobDescription(
		r.Context(),
		userID,
		*body.JobDescription,
	)
	if err != nil {
		a.Logger.Error("Failed to tailor profile: %v", err)
		http.Error(w, "Failed to tailor profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to generated response
	resp := a.toGeneratedJobMatch(jobMatch)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		a.Logger.Error("Failed to encode response: %v", err)
	}
}

func (a *API) toGeneratedJobMatch(jobMatch *domain.JobMatch) generated.JobMatch {
	id := jobMatch.ID
	matchScore := float32(jobMatch.MatchScore)
	tailoredSummary := jobMatch.TailoredSummary

	return generated.JobMatch{
		Id:              &id,
		MatchScore:      &matchScore,
		TailoredSummary: &tailoredSummary,
	}
}
