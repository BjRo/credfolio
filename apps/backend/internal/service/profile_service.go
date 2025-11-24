package service

import (
	"context"
	"fmt"
	"time"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/credfolio/apps/backend/pkg/logger"
	"github.com/credfolio/apps/backend/pkg/pdf"
	"github.com/google/uuid"
)

// ProfileService handles profile generation and management
type ProfileService struct {
	profileRepo         repository.ProfileRepository
	workExpRepo         repository.WorkExperienceRepository
	credibilityRepo     repository.CredibilityHighlightRepository
	referenceLetterRepo repository.ReferenceLetterRepository
	llmProvider         LLMProvider
	pdfExtractor        pdf.ExtractorInterface
	logger              *logger.Logger
}

// NewProfileService creates a new profile service
func NewProfileService(
	profileRepo repository.ProfileRepository,
	workExpRepo repository.WorkExperienceRepository,
	credibilityRepo repository.CredibilityHighlightRepository,
	referenceLetterRepo repository.ReferenceLetterRepository,
	llmProvider LLMProvider,
	pdfExtractor pdf.ExtractorInterface,
	logger *logger.Logger,
) *ProfileService {
	return &ProfileService{
		profileRepo:         profileRepo,
		workExpRepo:         workExpRepo,
		credibilityRepo:     credibilityRepo,
		referenceLetterRepo: referenceLetterRepo,
		llmProvider:         llmProvider,
		pdfExtractor:        pdfExtractor,
		logger:              logger,
	}
}

// GenerateProfileFromReferences generates a profile from uploaded reference letters
func (s *ProfileService) GenerateProfileFromReferences(ctx context.Context, userID uuid.UUID, referenceLetterIDs []uuid.UUID) (*domain.Profile, error) {
	s.logger.Info("Starting profile generation for user %s with %d reference letters", userID, len(referenceLetterIDs))

	// Get or create profile
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Profile doesn't exist, create a new one
		profile = &domain.Profile{
			UserID: userID,
		}
		if err := s.profileRepo.Create(ctx, profile); err != nil {
			return nil, fmt.Errorf("failed to create profile: %w", err)
		}
	}

	// Process each reference letter
	for _, letterID := range referenceLetterIDs {
		letter, err := s.referenceLetterRepo.GetByID(ctx, letterID)
		if err != nil {
			s.logger.Error("Failed to get reference letter %s: %v", letterID, err)
			continue
		}

		// Extract text from PDF if not already extracted
		text := letter.ExtractedText
		if text == "" {
			// TODO: Extract from PDF file using pdfExtractor
			// For now, we assume text is already extracted
			s.logger.Error("Reference letter %s has no extracted text", letterID)
			continue
		}

		// Extract structured data using LLM
		profileData, err := s.llmProvider.ExtractProfileData(ctx, text)
		if err != nil {
			s.logger.Error("Failed to extract profile data from letter %s: %v", letterID, err)
			continue
		}

		// Extract credibility highlights
		credibilityData, err := s.llmProvider.ExtractCredibility(ctx, text)
		if err != nil {
			s.logger.Error("Failed to extract credibility from letter %s: %v", letterID, err)
			// Continue even if credibility extraction fails
		}

		// Parse dates
		startDate, err := time.Parse("2006-01-02", profileData.StartDate)
		if err != nil {
			s.logger.Error("Failed to parse start date: %v", err)
			continue
		}

		var endDate *time.Time
		if profileData.EndDate != "" {
			parsed, err := time.Parse("2006-01-02", profileData.EndDate)
			if err == nil {
				endDate = &parsed
			}
		}

		// Create work experience
		workExp := &domain.WorkExperience{
			ProfileID:         profile.ID,
			CompanyName:       profileData.CompanyName,
			Role:              profileData.Role,
			StartDate:         startDate,
			EndDate:           endDate,
			Description:       profileData.Description,
			ReferenceLetterID: &letterID,
		}

		if err := s.workExpRepo.Create(ctx, workExp); err != nil {
			s.logger.Error("Failed to create work experience: %v", err)
			continue
		}

		// Create credibility highlights
		if credibilityData != nil {
			for _, quote := range credibilityData.Quotes {
				highlight := &domain.CredibilityHighlight{
					WorkExperienceID: workExp.ID,
					Quote:            quote,
					Sentiment:        domain.Sentiment(credibilityData.Sentiment),
					SourceLetterID:   letterID,
				}
				if err := s.credibilityRepo.Create(ctx, highlight); err != nil {
					s.logger.Error("Failed to create credibility highlight: %v", err)
				}
			}
		}

		// TODO: Create/update skills
		// This would require a SkillRepository and ProfileSkill management
	}

	// Reload profile with all associations
	profile, err = s.profileRepo.GetByID(ctx, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload profile: %w", err)
	}

	s.logger.Info("Profile generation completed for user %s", userID)
	return profile, nil
}

// GetProfile retrieves a profile by user ID
func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	return profile, nil
}

// UpdateProfile updates an existing profile
func (s *ProfileService) UpdateProfile(ctx context.Context, profile *domain.Profile) error {
	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}
