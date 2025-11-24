package handler

import (
	"encoding/json"
	"net/http"

	"github.com/credfolio/apps/backend/api/generated"
	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// GenerateProfile implements generated.ServerInterface
func (a *API) GenerateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == uuid.Nil {
		writeErrorResponse(w, http.StatusUnauthorized, ErrorCodeUnauthorized, "Unauthorized")
		return
	}

	// Get all reference letters for the user
	letters, err := a.ReferenceLetterRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		a.Logger.Error("Failed to get reference letters: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeDatabaseError, "Internal server error")
		return
	}

	if len(letters) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeNoReferenceLetters, "No reference letters found")
		return
	}

	var letterIDs []uuid.UUID
	for _, letter := range letters {
		letterIDs = append(letterIDs, letter.ID)
	}

	profile, err := a.ProfileService.GenerateProfileFromReferences(r.Context(), userID, letterIDs)
	if err != nil {
		a.Logger.Error("Failed to generate profile: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeProfileGenerationFailed, "Failed to generate profile")
		return
	}

	resp := a.toGeneratedProfile(profile)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		a.Logger.Error("Failed to encode response: %v", err)
	}
}

// GetProfile implements generated.ServerInterface
func (a *API) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == uuid.Nil {
		writeErrorResponse(w, http.StatusUnauthorized, ErrorCodeUnauthorized, "Unauthorized")
		return
	}

	profile, err := a.ProfileService.GetProfile(r.Context(), userID)
	if err != nil {
		a.Logger.Error("Failed to get profile: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeDatabaseError, "Internal server error")
		return
	}

	// If profile doesn't exist (might return nil or error depending on repo), handle it.
	if profile == nil {
		writeErrorResponse(w, http.StatusNotFound, ErrorCodeProfileNotFound, "Profile not found")
		return
	}

	resp := a.toGeneratedProfile(profile)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		a.Logger.Error("Failed to encode response: %v", err)
	}
}

// UpdateProfile implements generated.ServerInterface
func (a *API) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == uuid.Nil {
		writeErrorResponse(w, http.StatusUnauthorized, ErrorCodeUnauthorized, "Unauthorized")
		return
	}

	var input generated.ProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeInvalidRequestBody, "Invalid request body")
		return
	}

	// Validate profile summary
	if err := ValidateProfileSummary(input.Summary); err != nil {
		if valErr, ok := err.(*ValidationError); ok {
			writeErrorResponse(w, http.StatusBadRequest, valErr.ErrorCode, valErr.Message)
			return
		}
		writeErrorResponse(w, http.StatusBadRequest, ErrorCodeInvalidRequest, "Invalid profile summary")
		return
	}

	// Get existing profile
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

	// Update fields with sanitized values
	if input.Summary != nil {
		profile.Summary = SanitizeProfileSummary(*input.Summary)
	}

	// TODO: Update Work Experiences and Skills
	// For now, we only support updating Summary as MVP for editing.

	if err := a.ProfileService.UpdateProfile(r.Context(), profile); err != nil {
		a.Logger.Error("Failed to update profile: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, ErrorCodeProfileUpdateFailed, "Failed to update profile")
		return
	}

	resp := a.toGeneratedProfile(profile)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		a.Logger.Error("Failed to encode response: %v", err)
	}
}

func (a *API) toGeneratedProfile(profile *domain.Profile) generated.Profile {
	resp := generated.Profile{
		Id:      &profile.ID,
		Summary: &profile.Summary,
	}

	// Map WorkExperiences
	var workExps []generated.WorkExperience
	for _, we := range profile.WorkExperiences {
		id := we.ID
		start := openapi_types.Date{Time: we.StartDate}

		var end *openapi_types.Date
		if we.EndDate != nil {
			e := openapi_types.Date{Time: *we.EndDate}
			end = &e
		}

		// Map Credibility Highlights
		var highlights []generated.CredibilityHighlight
		for _, ch := range we.CredibilityHighlights {
			quote := ch.Quote
			sentiment := generated.CredibilityHighlightSentiment(ch.Sentiment)
			highlights = append(highlights, generated.CredibilityHighlight{
				Quote:     &quote,
				Sentiment: &sentiment,
			})
		}

		weResp := generated.WorkExperience{
			Id:                    &id,
			CompanyName:           &we.CompanyName,
			Role:                  &we.Role,
			StartDate:             &start,
			EndDate:               end,
			Description:           &we.Description,
			CredibilityHighlights: &highlights,
		}
		workExps = append(workExps, weResp)
	}
	resp.WorkExperiences = &workExps

	// Map Skills
	var skills []string
	for _, s := range profile.Skills {
		skills = append(skills, s.Name)
	}
	resp.Skills = &skills

	return resp
}
