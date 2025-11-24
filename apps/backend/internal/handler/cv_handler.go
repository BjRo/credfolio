package handler

import (
	"fmt"
	"net/http"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/credfolio/apps/backend/pkg/pdf"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// DownloadCV implements generated.ServerInterface
func (a *API) DownloadCV(w http.ResponseWriter, r *http.Request, profileId openapi_types.UUID) {
	userID := middleware.GetUserID(r)
	if userID == uuid.Nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert OpenAPI UUID to uuid.UUID
	profileUUID := uuid.UUID(profileId)

	// Get profile
	profile, err := a.ProfileService.GetProfile(r.Context(), userID)
	if err != nil {
		a.Logger.Error("Failed to get profile: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if profile == nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Verify profile belongs to user
	if profile.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Verify profile ID matches
	if profile.ID != profileUUID {
		http.Error(w, "Profile ID mismatch", http.StatusBadRequest)
		return
	}

	// Check for jobMatchId query parameter for tailored CV
	var jobMatch *domain.JobMatch
	jobMatchIdStr := r.URL.Query().Get("jobMatchId")
	if jobMatchIdStr != "" {
		jobMatchUUID, err := uuid.Parse(jobMatchIdStr)
		if err != nil {
			http.Error(w, "Invalid jobMatchId", http.StatusBadRequest)
			return
		}

		// Get job match and verify it belongs to this profile
		jobMatch, err = a.JobMatchRepo.GetByID(r.Context(), jobMatchUUID)
		if err != nil {
			a.Logger.Error("Failed to get job match: %v", err)
			http.Error(w, "Job match not found", http.StatusNotFound)
			return
		}

		// Verify job match belongs to this profile
		if jobMatch.BaseProfileID != profileUUID {
			http.Error(w, "Job match does not belong to this profile", http.StatusForbidden)
			return
		}
	}

	// Generate PDF
	cvGenerator := pdf.NewCVGenerator()
	pdfBytes, err := cvGenerator.GenerateCVFromProfile(profile, jobMatch)
	if err != nil {
		a.Logger.Error("Failed to generate CV: %v", err)
		http.Error(w, "Failed to generate CV", http.StatusInternalServerError)
		return
	}

	// Set headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=cv.pdf")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	// Write PDF
	if _, err := w.Write(pdfBytes); err != nil {
		a.Logger.Error("Failed to write PDF response: %v", err)
	}
}
