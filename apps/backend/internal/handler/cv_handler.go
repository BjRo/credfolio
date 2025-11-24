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

	profileUUID := uuid.UUID(profileId)

	profile, err := a.ProfileService.GetProfile(r.Context(), userID)
	if err != nil {
		a.Logger.Error("Failed to get profile: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeDatabaseError, "Internal server error")
		return
	}

	if !a.verifyProfile(w, profile, userID, profileUUID) {
		return
	}

	var jobMatch *domain.JobMatch
	jobMatchIdStr := r.URL.Query().Get("jobMatchId")
	if jobMatchIdStr != "" {
		jobMatchUUID, err := ValidateUUID(jobMatchIdStr, "jobMatchId")
		if err != nil {
			if valErr, ok := err.(*ValidationError); ok {
				writeErrorResponse(w, http.StatusBadRequest, valErr.ErrorCode, valErr.Message)
				return
			}
			writeErrorResponse(w, http.StatusBadRequest, ErrorCodeInvalidJobMatchID, "Invalid jobMatchId")
			return
		}

		jobMatch, err = a.JobMatchRepo.GetByID(r.Context(), jobMatchUUID)
		if err != nil {
			a.Logger.Error("Failed to get job match: %v", err)
			writeErrorResponse(w, http.StatusNotFound, ErrorCodeJobMatchNotFound, "Job match not found")
			return
		}

		if !a.verifyJobMatch(w, jobMatch, profileUUID) {
			return
		}
	}

	cvGenerator := pdf.NewCVGenerator()
	pdfBytes, err := cvGenerator.GenerateCVFromProfile(profile, jobMatch)
	if err != nil {
		a.Logger.Error("Failed to generate CV: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodePDFGenerationFailed, "Failed to generate CV")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=cv.pdf")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	if _, err := w.Write(pdfBytes); err != nil {
		a.Logger.Error("Failed to write PDF response: %v", err)
	}
}

// verifyProfile verifies that a profile exists, belongs to the user, and matches the expected profile ID.
// Returns false if verification fails (and writes error response), true otherwise.
func (a *API) verifyProfile(w http.ResponseWriter, profile *domain.Profile, userID uuid.UUID, profileUUID uuid.UUID) bool {
	if profile == nil {
		writeErrorResponse(w, http.StatusNotFound, ErrorCodeProfileNotFound, "Profile not found")
		return false
	}

	if profile.UserID != userID {
		writeErrorResponse(w, http.StatusForbidden, ErrorCodeForbidden, "Forbidden")
		return false
	}

	if profile.ID != profileUUID {
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeProfileIDMismatch, "Profile ID mismatch")
		return false
	}

	return true
}

// verifyJobMatch verifies that a job match belongs to the specified profile.
// Returns false if verification fails (and writes error response), true otherwise.
func (a *API) verifyJobMatch(w http.ResponseWriter, jobMatch *domain.JobMatch, profileUUID uuid.UUID) bool {
	if jobMatch.BaseProfileID != profileUUID {
		writeErrorResponse(w, http.StatusForbidden, ErrorCodeJobMatchMismatch, "Job match does not belong to this profile")
		return false
	}
	return true
}
