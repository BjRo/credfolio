package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/credfolio/apps/backend/pkg/logger"
	"github.com/credfolio/apps/backend/pkg/pdf"
	"github.com/google/uuid"
)

// ProfileServicer defines the interface for profile service operations
type ProfileServicer interface {
	GenerateProfileFromReferences(ctx context.Context, userID uuid.UUID) (*domain.Profile, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, profile *domain.Profile) error
}

// TailoringServicer defines the interface for tailoring service operations
type TailoringServicer interface {
	TailorProfileToJobDescription(ctx context.Context, userID uuid.UUID, jobDescription string) (*TailoredProfileResponse, error)
}

// ProfileHandler handles profile HTTP requests
type ProfileHandler struct {
	profileService   ProfileServicer
	tailoringService TailoringServicer
	logger           *logger.Logger
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(profileService ProfileServicer) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		logger:         logger.New(logger.LevelInfo),
	}
}

// NewProfileHandlerWithTailoring creates a new profile handler with tailoring support
func NewProfileHandlerWithTailoring(profileService ProfileServicer, tailoringService TailoringServicer) *ProfileHandler {
	return &ProfileHandler{
		profileService:   profileService,
		tailoringService: tailoringService,
		logger:           logger.New(logger.LevelInfo),
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

// TailoredProfileResponse represents the API response for a tailored profile
type TailoredProfileResponse struct {
	ID                  string                       `json:"id"`
	MatchScore          float64                      `json:"matchScore"`
	MatchSummary        string                       `json:"matchSummary"`
	TailoredExperiences []TailoredExperienceResponse `json:"tailoredExperiences"`
	RelevantSkills      []string                     `json:"relevantSkills"`
}

// TailoredExperienceResponse represents a tailored work experience in the API response
type TailoredExperienceResponse struct {
	ID              string  `json:"id"`
	CompanyName     string  `json:"companyName"`
	Role            string  `json:"role"`
	StartDate       string  `json:"startDate"`
	EndDate         string  `json:"endDate,omitempty"`
	Description     string  `json:"description"`
	RelevanceScore  float64 `json:"relevanceScore"`
	HighlightReason string  `json:"highlightReason,omitempty"`
}

// TailorRequest represents the request body for tailoring a profile
type TailorRequest struct {
	JobDescription string `json:"jobDescription"`
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

// DownloadCV handles GET /profile/cv - download profile as PDF CV
func (h *ProfileHandler) DownloadCV(w http.ResponseWriter, r *http.Request) {
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

	// Convert domain profile to PDF data
	pdfData := h.toProfilePDFData(profile)

	// Generate PDF
	generator := pdf.NewGenerator()
	pdfBytes, err := generator.GenerateCV(pdfData)
	if err != nil {
		h.logger.Error("Failed to generate CV: %v", err)
		http.Error(w, "Failed to generate CV", http.StatusInternalServerError)
		return
	}

	// Send PDF response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"cv-%s.pdf\"", userID.String()[:8]))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	_, _ = w.Write(pdfBytes)
}

// toProfilePDFData converts a domain profile to PDF data
func (h *ProfileHandler) toProfilePDFData(profile *domain.Profile) *pdf.ProfileData {
	data := &pdf.ProfileData{
		Name:    "Professional", // Default name, would come from user in real app
		Summary: profile.Summary,
		Skills:  make([]string, 0),
	}

	for _, we := range profile.WorkExperiences {
		expData := pdf.WorkExperienceData{
			CompanyName: we.CompanyName,
			Role:        we.Role,
			StartDate:   we.StartDate.Format("2006-01"),
			Description: we.Description,
		}
		if we.EndDate != nil {
			expData.EndDate = we.EndDate.Format("2006-01")
		}

		for _, ch := range we.CredibilityHighlights {
			expData.Highlights = append(expData.Highlights, pdf.HighlightData{
				Quote:     ch.Quote,
				Sentiment: string(ch.Sentiment),
			})
		}

		data.WorkExperiences = append(data.WorkExperiences, expData)
	}

	for _, skill := range profile.Skills {
		data.Skills = append(data.Skills, skill.Name)
	}

	return data
}

// Tailor handles POST /profile/tailor - tailor profile to job description
func (h *ProfileHandler) Tailor(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input TailorRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.JobDescription == "" {
		http.Error(w, "Job description is required", http.StatusBadRequest)
		return
	}

	if h.tailoringService == nil {
		http.Error(w, "Tailoring service not available", http.StatusServiceUnavailable)
		return
	}

	h.logger.Info("Tailoring profile for user %s", userID.String())

	tailoredProfile, err := h.tailoringService.TailorProfileToJobDescription(r.Context(), userID, input.JobDescription)
	if err != nil {
		h.logger.Error("Failed to tailor profile: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tailoredProfile)
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
