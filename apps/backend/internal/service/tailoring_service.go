package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/credfolio/apps/backend/pkg/logger"
	"github.com/google/uuid"
)

// TailoringService handles profile tailoring to job descriptions
type TailoringService struct {
	profileRepo   repository.ProfileRepository
	jobMatchRepo  repository.JobMatchRepository
	llmProvider   LLMProvider
	logger        *logger.Logger
}

// NewTailoringService creates a new tailoring service
func NewTailoringService(
	profileRepo repository.ProfileRepository,
	jobMatchRepo repository.JobMatchRepository,
	llmProvider LLMProvider,
	logger *logger.Logger,
) *TailoringService {
	return &TailoringService{
		profileRepo:  profileRepo,
		jobMatchRepo: jobMatchRepo,
		llmProvider:  llmProvider,
		logger:       logger,
	}
}

// TailorProfileToJobDescription creates a tailored profile for a job description
func (s *TailoringService) TailorProfileToJobDescription(
	ctx context.Context,
	userID uuid.UUID,
	jobDescription string,
) (*domain.JobMatch, error) {
	s.logger.Info("Tailoring profile for user %s", userID)

	// Validate job description
	if strings.TrimSpace(jobDescription) == "" {
		return nil, fmt.Errorf("job description cannot be empty")
	}

	// Get user's profile
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Convert profile to text format for LLM
	profileText := s.buildProfileText(profile)

	// Generate tailored summary and match score using LLM
	tailoredSummary, matchScore, err := s.llmProvider.TailorProfile(ctx, profileText, jobDescription)
	if err != nil {
		return nil, fmt.Errorf("failed to tailor profile: %w", err)
	}

	// Create JobMatch record
	jobMatch := &domain.JobMatch{
		BaseProfileID:   profile.ID,
		JobDescription:  jobDescription,
		MatchScore:      matchScore,
		TailoredSummary: tailoredSummary,
	}

	if err := s.jobMatchRepo.Create(ctx, jobMatch); err != nil {
		return nil, fmt.Errorf("failed to create job match: %w", err)
	}

	s.logger.Info("Created job match %s with score %.2f", jobMatch.ID, matchScore)
	return jobMatch, nil
}

// buildProfileText converts a profile to a text representation for LLM processing
func (s *TailoringService) buildProfileText(profile *domain.Profile) string {
	var parts []string

	if profile.Summary != "" {
		parts = append(parts, fmt.Sprintf("Summary: %s", profile.Summary))
	}

	// Add work experiences
	if len(profile.WorkExperiences) > 0 {
		parts = append(parts, "\nWork Experience:")
		for _, exp := range profile.WorkExperiences {
			expText := fmt.Sprintf("- %s at %s", exp.Role, exp.CompanyName)
			if exp.Description != "" {
				expText += fmt.Sprintf("\n  %s", exp.Description)
			}
			if len(exp.CredibilityHighlights) > 0 {
				expText += "\n  Credibility Highlights:"
				for _, highlight := range exp.CredibilityHighlights {
					expText += fmt.Sprintf("\n    \"%s\" (%s)", highlight.Quote, highlight.Sentiment)
				}
			}
			parts = append(parts, expText)
		}
	}

	// Add skills
	if len(profile.Skills) > 0 {
		skillNames := make([]string, len(profile.Skills))
		for i, skill := range profile.Skills {
			skillNames[i] = skill.Name
		}
		parts = append(parts, fmt.Sprintf("\nSkills: %s", strings.Join(skillNames, ", ")))
	}

	return strings.Join(parts, "\n")
}
