package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/credfolio/apps/backend/internal/service"
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

// TailoredProfileResponse is defined in service package, re-exported here for convenience
type TailoredProfileResponse = service.TailoredProfileResponse
type TailoredExperienceResponse = service.TailoredExperienceResponse

// TailoringServicer defines the interface for tailoring service operations
type TailoringServicer interface {
	TailorProfileToJobDescription(ctx context.Context, userID uuid.UUID, jobDescription string) (*service.TailoredProfileResponse, error)
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

// TailorRequest represents the request body for tailoring a profile
type TailorRequest struct {
	JobDescription string `json:"jobDescription"`
}

const (
	// MinJobDescriptionLength is the minimum allowed job description length
	MinJobDescriptionLength = 50
	// MaxJobDescriptionLength is the maximum allowed job description length
	MaxJobDescriptionLength = 10000
	// MaxSummaryLength is the maximum allowed profile summary length
	MaxSummaryLength = 2000
)

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
		h.logger.Error("User ID not found in context for tailor request")
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input TailorRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("Failed to decode tailor request body: %v", err)
		h.writeError(w, "Invalid request body. Please provide a valid JSON object.", http.StatusBadRequest)
		return
	}

	// Validate job description
	jobDesc := strings.TrimSpace(input.JobDescription)
	if jobDesc == "" {
		h.writeError(w, "Job description is required", http.StatusBadRequest)
		return
	}
	if len(jobDesc) < MinJobDescriptionLength {
		h.writeError(w, fmt.Sprintf("Job description too short. Minimum %d characters required for accurate matching.", MinJobDescriptionLength), http.StatusBadRequest)
		return
	}
	if len(jobDesc) > MaxJobDescriptionLength {
		h.writeError(w, fmt.Sprintf("Job description too long. Maximum %d characters allowed.", MaxJobDescriptionLength), http.StatusBadRequest)
		return
	}

	if h.tailoringService == nil {
		h.logger.Error("Tailoring service not configured")
		h.writeError(w, "Profile tailoring is not available at this time", http.StatusServiceUnavailable)
		return
	}

	h.logger.Info("Tailoring profile for user %s (job description length: %d chars)", userID.String(), len(jobDesc))

	tailoredProfile, err := h.tailoringService.TailorProfileToJobDescription(r.Context(), userID, jobDesc)
	if err != nil {
		h.logger.Error("Failed to tailor profile for user %s: %v", userID.String(), err)
		// Check for specific error types
		if strings.Contains(err.Error(), "not found") {
			h.writeError(w, "Profile not found. Please generate your profile first.", http.StatusNotFound)
			return
		}
		h.writeError(w, "Failed to tailor profile. Please try again.", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Profile tailored successfully for user %s (match score: %.1f%%)", userID.String(), tailoredProfile.MatchScore)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tailoredProfile)
}

// writeError writes a JSON error response
func (h *ProfileHandler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
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
