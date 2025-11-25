package handler

import (
	"encoding/json"
	"net/http"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/credfolio/apps/backend/internal/service"
	"github.com/credfolio/apps/backend/pkg/logger"
)

// ProfileHandler handles profile HTTP requests
type ProfileHandler struct {
	profileService *service.ProfileService
	logger         *logger.Logger
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		logger:         logger.New(logger.LevelInfo),
	}
}

// ProfileResponse represents the API response for a profile
type ProfileResponse struct {
	ID              string                   `json:"id"`
	Summary         string                   `json:"summary"`
	WorkExperiences []WorkExperienceResponse `json:"workExperiences"`
	Skills          []string                 `json:"skills"`
}

// WorkExperienceResponse represents a work experience in the API response
type WorkExperienceResponse struct {
	ID                    string                         `json:"id"`
	CompanyName           string                         `json:"companyName"`
	Role                  string                         `json:"role"`
	StartDate             string                         `json:"startDate"`
	EndDate               string                         `json:"endDate,omitempty"`
	Description           string                         `json:"description"`
	CredibilityHighlights []CredibilityHighlightResponse `json:"credibilityHighlights,omitempty"`
}

// CredibilityHighlightResponse represents a credibility highlight in the API response
type CredibilityHighlightResponse struct {
	Quote     string `json:"quote"`
	Sentiment string `json:"sentiment"`
}

// Get handles GET /profile - get current user's profile
func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get profile: %v", err)
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	if profile == nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	response := h.toProfileResponse(profile)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// Generate handles POST /profile/generate - generate profile from reference letters
func (h *ProfileHandler) Generate(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.logger.Info("Starting profile generation for user %s", userID.String())

	profile, err := h.profileService.GenerateProfileFromReferences(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to generate profile: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := h.toProfileResponse(profile)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// Update handles PUT /profile - update profile
func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		Summary string `json:"summary"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	profile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get profile: %v", err)
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	if profile == nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	profile.Summary = input.Summary

	if err := h.profileService.UpdateProfile(r.Context(), profile); err != nil {
		h.logger.Error("Failed to update profile: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	response := h.toProfileResponse(profile)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// toProfileResponse converts a domain profile to an API response
func (h *ProfileHandler) toProfileResponse(profile *domain.Profile) *ProfileResponse {
	response := &ProfileResponse{
		ID:              profile.ID.String(),
		Summary:         profile.Summary,
		WorkExperiences: make([]WorkExperienceResponse, 0),
		Skills:          make([]string, 0),
	}

	for _, we := range profile.WorkExperiences {
		weResponse := WorkExperienceResponse{
			ID:          we.ID.String(),
			CompanyName: we.CompanyName,
			Role:        we.Role,
			StartDate:   we.StartDate.Format("2006-01-02"),
			Description: we.Description,
		}

		if we.EndDate != nil {
			weResponse.EndDate = we.EndDate.Format("2006-01-02")
		}

		for _, ch := range we.CredibilityHighlights {
			weResponse.CredibilityHighlights = append(weResponse.CredibilityHighlights, CredibilityHighlightResponse{
				Quote:     ch.Quote,
				Sentiment: string(ch.Sentiment),
			})
		}

		response.WorkExperiences = append(response.WorkExperiences, weResponse)
	}

	for _, skill := range profile.Skills {
		response.Skills = append(response.Skills, skill.Name)
	}

	return response
}
