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
		writeErrorResponse(w, http.StatusUnauthorized, ErrorCodeUnauthorized, "Unauthorized")
		return
	}

	// Convert OpenAPI UUID to uuid.UUID
	profileUUID := uuid.UUID(profileId)

	// Get profile
	profile, err := a.ProfileService.GetProfile(r.Context(), userID)
	if err != nil {
		a.Logger.Error("Failed to get profile: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeDatabaseError, "Internal server error")
		return
	}

	if profile == nil {
		writeErrorResponse(w, http.StatusNotFound, ErrorCodeProfileNotFound, "Profile not found")
		return
	}

	// Verify profile belongs to user
	if profile.UserID != userID {
		writeErrorResponse(w, http.StatusForbidden, ErrorCodeForbidden, "Forbidden")
		return
	}

	// Verify profile ID matches
	if profile.ID != profileUUID {
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeProfileIDMismatch, "Profile ID mismatch")
		return
	}

	// Check for jobMatchId query parameter for tailored CV
	var jobMatch *domain.JobMatch
	jobMatchIdStr := r.URL.Query().Get("jobMatchId")
	if jobMatchIdStr != "" {
		jobMatchUUID, err := uuid.Parse(jobMatchIdStr)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, ErrorCodeInvalidJobMatchID, "Invalid jobMatchId")
			return
		}

		// Get job match and verify it belongs to this profile
		jobMatch, err = a.JobMatchRepo.GetByID(r.Context(), jobMatchUUID)
		if err != nil {
			a.Logger.Error("Failed to get job match: %v", err)
			writeErrorResponse(w, http.StatusNotFound, ErrorCodeJobMatchNotFound, "Job match not found")
			return
		}

		// Verify job match belongs to this profile
		if jobMatch.BaseProfileID != profileUUID {
			writeErrorResponse(w, http.StatusForbidden, ErrorCodeJobMatchMismatch, "Job match does not belong to this profile")
			return
		}
	}

	// Generate PDF
	cvGenerator := pdf.NewCVGenerator()
	pdfBytes, err := cvGenerator.GenerateCVFromProfile(profile, jobMatch)
	if err != nil {
		a.Logger.Error("Failed to generate CV: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodePDFGenerationFailed, "Failed to generate CV")
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
