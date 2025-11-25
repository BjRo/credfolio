package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/google/uuid"
)

// TailoringService handles profile tailoring to job descriptions
type TailoringService struct {
	profileRepo  repository.ProfileRepository
	jobMatchRepo repository.JobMatchRepository
	llm          LLMProvider
}

// TailoredExperience represents a work experience with its relevance score
type TailoredExperience struct {
	Experience      domain.WorkExperience `json:"experience"`
	RelevanceScore  float64               `json:"relevanceScore"`
	HighlightReason string                `json:"highlightReason"`
}

// TailoredProfile represents a profile tailored to a specific job description
type TailoredProfile struct {
	ID                  uuid.UUID            `json:"id"`
	BaseProfile         *domain.Profile      `json:"baseProfile"`
	JobDescription      string               `json:"jobDescription"`
	MatchScore          float64              `json:"matchScore"`
	MatchSummary        string               `json:"matchSummary"`
	TailoredExperiences []TailoredExperience `json:"tailoredExperiences"`
	RelevantSkills      []string             `json:"relevantSkills"`
	CreatedAt           time.Time            `json:"createdAt"`
}

// NewTailoringService creates a new tailoring service
func NewTailoringService(
	profileRepo repository.ProfileRepository,
	jobMatchRepo repository.JobMatchRepository,
	llm LLMProvider,
) *TailoringService {
	return &TailoringService{
		profileRepo:  profileRepo,
		jobMatchRepo: jobMatchRepo,
		llm:          llm,
	}
}

// TailorProfileToJobDescription tailors a user's profile to a specific job description
func (s *TailoringService) TailorProfileToJobDescription(
	ctx context.Context,
	userID uuid.UUID,
	jobDescription string,
) (*TailoredProfile, error) {
	// Validate job description
	if jobDescription == "" {
		return nil, errors.New("job description is required")
	}

	// Get the user's profile
	profile, err := s.profileRepo.FindByUserIDWithRelations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find profile: %w", err)
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}

	// Calculate relevance scores for each experience
	tailoredExperiences := make([]TailoredExperience, 0, len(profile.WorkExperiences))
	var totalRelevance float64

	for _, exp := range profile.WorkExperiences {
		// Create a text representation of the experience
		expText := fmt.Sprintf("%s at %s: %s", exp.Role, exp.CompanyName, exp.Description)

		// Calculate relevance score using LLM
		relevance, err := s.llm.CalculateRelevance(ctx, expText, jobDescription)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate relevance: %w", err)
		}

		tailoredExperiences = append(tailoredExperiences, TailoredExperience{
			Experience:     exp,
			RelevanceScore: relevance,
		})
		totalRelevance += relevance
	}

	// Sort experiences by relevance (highest first)
	sort.Slice(tailoredExperiences, func(i, j int) bool {
		return tailoredExperiences[i].RelevanceScore > tailoredExperiences[j].RelevanceScore
	})

	// Calculate overall match score
	var matchScore float64
	if len(tailoredExperiences) > 0 {
		matchScore = totalRelevance / float64(len(tailoredExperiences))
	}

	// Ensure score is between 0 and 1
	if matchScore < 0 {
		matchScore = 0
	}
	if matchScore > 1 {
		matchScore = 1
	}

	// Generate match summary
	summaryPrompt := fmt.Sprintf(
		"Based on the profile summary '%s' and job description '%s', "+
			"write a brief 1-2 sentence explanation of why this candidate matches or doesn't match the job.",
		profile.Summary,
		jobDescription,
	)
	matchSummary, err := s.llm.GenerateText(ctx, summaryPrompt)
	if err != nil {
		matchSummary = "Match analysis unavailable"
	}

	// Extract relevant skills from profile
	relevantSkills := make([]string, 0)
	for _, skill := range profile.Skills {
		relevantSkills = append(relevantSkills, skill.Name)
	}

	// Create job match record
	jobMatch := &domain.JobMatch{
		ID:              uuid.New(),
		BaseProfileID:   profile.ID,
		JobDescription:  jobDescription,
		MatchScore:      matchScore,
		TailoredSummary: matchSummary,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.jobMatchRepo.Create(ctx, jobMatch); err != nil {
		return nil, fmt.Errorf("failed to save job match: %w", err)
	}

	return &TailoredProfile{
		ID:                  jobMatch.ID,
		BaseProfile:         profile,
		JobDescription:      jobDescription,
		MatchScore:          matchScore,
		MatchSummary:        matchSummary,
		TailoredExperiences: tailoredExperiences,
		RelevantSkills:      relevantSkills,
		CreatedAt:           jobMatch.CreatedAt,
	}, nil
}

// GetTailoredProfile retrieves a previously tailored profile by ID
func (s *TailoringService) GetTailoredProfile(ctx context.Context, id uuid.UUID) (*domain.JobMatch, error) {
	return s.jobMatchRepo.FindByID(ctx, id)
}

// GetTailoredProfilesByUser retrieves all tailored profiles for a user
func (s *TailoringService) GetTailoredProfilesByUser(ctx context.Context, userID uuid.UUID) ([]*domain.JobMatch, error) {
	profile, err := s.profileRepo.FindByUserIDWithRelations(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, nil
	}

	return s.jobMatchRepo.FindByProfileID(ctx, profile.ID)
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

// TailoredProfileResponse represents the API response for a tailored profile
type TailoredProfileResponse struct {
	ID                  string                       `json:"id"`
	MatchScore          float64                      `json:"matchScore"`
	MatchSummary        string                       `json:"matchSummary"`
	TailoredExperiences []TailoredExperienceResponse `json:"tailoredExperiences"`
	RelevantSkills      []string                     `json:"relevantSkills"`
}

// TailoringServiceAdapter wraps TailoringService to convert response types
type TailoringServiceAdapter struct {
	*TailoringService
}

// NewTailoringServiceAdapter creates a new adapter
func NewTailoringServiceAdapter(
	profileRepo repository.ProfileRepository,
	jobMatchRepo repository.JobMatchRepository,
	llm LLMProvider,
) *TailoringServiceAdapter {
	return &TailoringServiceAdapter{
		TailoringService: NewTailoringService(profileRepo, jobMatchRepo, llm),
	}
}

// TailorProfileToJobDescription adapts the response to the handler's expected type
func (a *TailoringServiceAdapter) TailorProfileToJobDescription(
	ctx context.Context,
	userID uuid.UUID,
	jobDescription string,
) (*TailoredProfileResponse, error) {
	result, err := a.TailoringService.TailorProfileToJobDescription(ctx, userID, jobDescription)
	if err != nil {
		return nil, err
	}

	// Convert to response type
	experiences := make([]TailoredExperienceResponse, 0, len(result.TailoredExperiences))
	for _, exp := range result.TailoredExperiences {
		expResp := TailoredExperienceResponse{
			ID:              exp.Experience.ID.String(),
			CompanyName:     exp.Experience.CompanyName,
			Role:            exp.Experience.Role,
			StartDate:       exp.Experience.StartDate.Format("2006-01-02"),
			Description:     exp.Experience.Description,
			RelevanceScore:  exp.RelevanceScore,
			HighlightReason: exp.HighlightReason,
		}
		if exp.Experience.EndDate != nil {
			expResp.EndDate = exp.Experience.EndDate.Format("2006-01-02")
		}
		experiences = append(experiences, expResp)
	}

	return &TailoredProfileResponse{
		ID:                  result.ID.String(),
		MatchScore:          result.MatchScore,
		MatchSummary:        result.MatchSummary,
		TailoredExperiences: experiences,
		RelevantSkills:      result.RelevantSkills,
	}, nil
}
