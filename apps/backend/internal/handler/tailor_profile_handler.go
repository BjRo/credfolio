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
		writeErrorResponse(w, http.StatusUnauthorized, ErrorCodeUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var body generated.TailorProfileJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeInvalidRequestBody, "Invalid request body")
		return
	}

	// Validate job description
	if err := ValidateJobDescription(body.JobDescription); err != nil {
		if valErr, ok := err.(*ValidationError); ok {
			writeErrorResponse(w, http.StatusBadRequest, valErr.ErrorCode, valErr.Message)
			return
		}
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeJobDescriptionRequired, "Job description is required")
		return
	}

	// Sanitize job description
	sanitizedJobDescription := SanitizeJobDescription(*body.JobDescription)

	// Call tailoring service
	jobMatch, err := a.TailoringService.TailorProfileToJobDescription(
		r.Context(),
		userID,
		sanitizedJobDescription,
	)
	if err != nil {
		a.Logger.Error("Failed to tailor profile: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeProfileTailoringFailed, "Failed to tailor profile")
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
