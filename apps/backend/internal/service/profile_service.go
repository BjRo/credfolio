package service

import (
	"context"
	"time"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/credfolio/apps/backend/pkg/logger"
	"github.com/credfolio/apps/backend/pkg/pdf"
	"github.com/google/uuid"
)

// ProfileService handles profile generation and management
type ProfileService struct {
	*BaseService
	profileRepo         repository.ProfileRepository
	workExpRepo         repository.WorkExperienceRepository
	credibilityRepo     repository.CredibilityHighlightRepository
	referenceLetterRepo repository.ReferenceLetterRepository
	llmProvider         LLMProvider
	pdfExtractor        pdf.ExtractorInterface
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
		BaseService:         NewBaseService(logger),
		profileRepo:         profileRepo,
		workExpRepo:         workExpRepo,
		credibilityRepo:     credibilityRepo,
		referenceLetterRepo: referenceLetterRepo,
		llmProvider:         llmProvider,
		pdfExtractor:        pdfExtractor,
	}
}

// GenerateProfileFromReferences generates a profile from uploaded reference letters
func (s *ProfileService) GenerateProfileFromReferences(ctx context.Context, userID uuid.UUID, referenceLetterIDs []uuid.UUID) (*domain.Profile, error) {
	s.LogOperationStart("Starting profile generation for user %s with %d reference letters", userID, len(referenceLetterIDs))

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Profile doesn't exist, create a new one
		profile = &domain.Profile{
			UserID: userID,
		}
		if err := s.profileRepo.Create(ctx, profile); err != nil {
			return nil, s.WrapError(err, "failed to create profile")
		}
	}

	// Process each reference letter
	for _, letterID := range referenceLetterIDs {
		letter, err := s.referenceLetterRepo.GetByID(ctx, letterID)
		if err != nil {
			s.LogError("Failed to get reference letter %s: %v", letterID, err)
			continue
		}

		text := letter.ExtractedText
		if text == "" {
			// TODO: Extract from PDF file using pdfExtractor
			// For now, we assume text is already extracted
			s.LogError("Reference letter %s has no extracted text", letterID)
			continue
		}

		profileData, err := s.llmProvider.ExtractProfileData(ctx, text)
		if err != nil {
			s.LogError("Failed to extract profile data from letter %s: %v", letterID, err)
			continue
		}

		credibilityData, err := s.llmProvider.ExtractCredibility(ctx, text)
		if err != nil {
			s.LogError("Failed to extract credibility from letter %s: %v", letterID, err)
			// Continue even if credibility extraction fails
		}

		startDate, err := time.Parse("2006-01-02", profileData.StartDate)
		if err != nil {
			s.LogError("Failed to parse start date '%s' from letter %s: %v", profileData.StartDate, letterID, err)
			continue
		}

		var endDate *time.Time
		if profileData.EndDate != "" {
			parsed, err := time.Parse("2006-01-02", profileData.EndDate)
			if err == nil {
				endDate = &parsed
			}
		}

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
			s.LogError("Failed to create work experience: %v", err)
			continue
		}

		if credibilityData != nil {
			for _, quote := range credibilityData.Quotes {
				highlight := &domain.CredibilityHighlight{
					WorkExperienceID: workExp.ID,
					Quote:            quote,
					Sentiment:        domain.Sentiment(credibilityData.Sentiment),
					SourceLetterID:   letterID,
				}
				if err := s.credibilityRepo.Create(ctx, highlight); err != nil {
					s.LogError("Failed to create credibility highlight: %v", err)
				}
			}
		}

		// TODO: Create/update skills
		// This would require a SkillRepository and ProfileSkill management
	}

	profile, err = s.profileRepo.GetByID(ctx, profile.ID)
	if err != nil {
		return nil, s.WrapError(err, "failed to reload profile")
	}

	s.LogOperationComplete("Profile generation completed for user %s", userID)
	return profile, nil
}

// GetProfile retrieves a profile by user ID with all associated data
// This includes:
// - WorkExperiences with CredibilityHighlights (aggregated from all reference letters)
// - Skills aggregated across all work experiences (ManyToMany relationship ensures deduplication)
// - JobMatches if any
func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, s.WrapError(err, "failed to get profile")
	}
	return profile, nil
}

// UpdateProfile updates an existing profile
func (s *ProfileService) UpdateProfile(ctx context.Context, profile *domain.Profile) error {
	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return s.WrapError(err, "failed to update profile")
	}
	return nil
}
