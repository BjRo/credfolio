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
	*BaseService
	profileRepo  repository.ProfileRepository
	jobMatchRepo repository.JobMatchRepository
	llmProvider  LLMProvider
}

// NewTailoringService creates a new tailoring service
func NewTailoringService(
	profileRepo repository.ProfileRepository,
	jobMatchRepo repository.JobMatchRepository,
	llmProvider LLMProvider,
	logger *logger.Logger,
) *TailoringService {
	return &TailoringService{
		BaseService:  NewBaseService(logger),
		profileRepo:  profileRepo,
		jobMatchRepo: jobMatchRepo,
		llmProvider:  llmProvider,
	}
}

// TailorProfileToJobDescription creates a tailored profile for a job description
func (s *TailoringService) TailorProfileToJobDescription(
	ctx context.Context,
	userID uuid.UUID,
	jobDescription string,
) (*domain.JobMatch, error) {
	s.LogOperationStart("Tailoring profile for user %s", userID)

	if err := s.ValidateNotEmpty(jobDescription, "job description"); err != nil {
		return nil, err
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, s.WrapError(err, "failed to get profile")
	}

	profileText := s.buildProfileText(profile)

	tailoredSummary, matchScore, err := s.llmProvider.TailorProfile(ctx, profileText, jobDescription)
	if err != nil {
		return nil, s.WrapError(err, "failed to tailor profile")
	}

	jobMatch := &domain.JobMatch{
		BaseProfileID:   profile.ID,
		JobDescription:  jobDescription,
		MatchScore:      matchScore,
		TailoredSummary: tailoredSummary,
	}

	if err := s.jobMatchRepo.Create(ctx, jobMatch); err != nil {
		return nil, s.WrapError(err, "failed to create job match")
	}

	s.LogOperationComplete("Created job match %s with score %.2f", jobMatch.ID, matchScore)
	return jobMatch, nil
}

// buildProfileText converts a profile to a text representation for LLM processing
func (s *TailoringService) buildProfileText(profile *domain.Profile) string {
	var parts []string

	if profile.Summary != "" {
		parts = append(parts, fmt.Sprintf("Summary: %s", profile.Summary))
	}

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

	if len(profile.Skills) > 0 {
		skillNames := make([]string, len(profile.Skills))
		for i, skill := range profile.Skills {
			skillNames[i] = skill.Name
		}
		parts = append(parts, fmt.Sprintf("\nSkills: %s", strings.Join(skillNames, ", ")))
	}

	return strings.Join(parts, "\n")
}
