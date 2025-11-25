package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/credfolio/apps/backend/pkg/logger"
	"github.com/google/uuid"
)

// RateLimiter provides simple per-user rate limiting for AI API calls
type RateLimiter struct {
	mu           sync.Mutex
	userRequests map[uuid.UUID][]time.Time
	maxRequests  int
	window       time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		userRequests: make(map[uuid.UUID][]time.Time),
		maxRequests:  maxRequests,
		window:       window,
	}
}

// Allow checks if a request is allowed for the given user
func (rl *RateLimiter) Allow(userID uuid.UUID) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Filter out old requests
	requests := rl.userRequests[userID]
	var recent []time.Time
	for _, t := range requests {
		if t.After(windowStart) {
			recent = append(recent, t)
		}
	}

	// Check if under limit
	if len(recent) >= rl.maxRequests {
		return false
	}

	// Record this request
	recent = append(recent, now)
	rl.userRequests[userID] = recent
	return true
}

// ProfileService handles profile-related business logic
type ProfileService struct {
	profileRepo         repository.ProfileRepository
	referenceLetterRepo repository.ReferenceLetterRepository
	workExperienceRepo  repository.WorkExperienceRepository
	highlightRepo       repository.CredibilityHighlightRepository
	llmProvider         LLMProvider
	rateLimiter         *RateLimiter
	logger              *logger.Logger
}

// NewProfileService creates a new profile service
func NewProfileService(
	profileRepo repository.ProfileRepository,
	referenceLetterRepo repository.ReferenceLetterRepository,
	workExperienceRepo repository.WorkExperienceRepository,
	highlightRepo repository.CredibilityHighlightRepository,
	llmProvider LLMProvider,
) *ProfileService {
	return &ProfileService{
		profileRepo:         profileRepo,
		referenceLetterRepo: referenceLetterRepo,
		workExperienceRepo:  workExperienceRepo,
		highlightRepo:       highlightRepo,
		llmProvider:         llmProvider,
		rateLimiter:         NewRateLimiter(10, time.Minute), // 10 AI calls per minute per user
		logger:              logger.New(logger.LevelInfo),
	}
}

// ErrRateLimitExceeded is returned when rate limit is exceeded
var ErrRateLimitExceeded = errors.New("rate limit exceeded: too many AI requests, please wait before trying again")

// GenerateProfileFromReferences generates a profile from uploaded reference letters
func (s *ProfileService) GenerateProfileFromReferences(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	s.logger.Info("Starting profile generation for user %s", userID.String())

	// Check rate limit before AI operations
	if !s.rateLimiter.Allow(userID) {
		s.logger.Warn("Rate limit exceeded for user %s", userID.String())
		return nil, ErrRateLimitExceeded
	}

	// Find all pending reference letters for the user
	letters, err := s.referenceLetterRepo.FindPendingByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to find reference letters: %v", err)
		return nil, err
	}

	if len(letters) == 0 {
		return nil, errors.New("no pending reference letters found")
	}

	// Get or create profile
	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to find profile: %v", err)
		return nil, err
	}

	if profile == nil {
		profile = domain.NewProfile(userID)
		if err := s.profileRepo.Create(ctx, profile); err != nil {
			s.logger.Error("Failed to create profile: %v", err)
			return nil, err
		}
	}

	// Process each reference letter
	for _, letter := range letters {
		if letter.ExtractedText == "" {
			s.logger.Warn("Reference letter %s has no extracted text, skipping", letter.ID.String())
			continue
		}

		s.logger.Info("Processing reference letter %s", letter.ID.String())

		// Extract profile data using LLM
		profileData, err := s.llmProvider.ExtractProfileData(ctx, letter.ExtractedText)
		if err != nil {
			s.logger.Error("Failed to extract profile data from letter %s: %v", letter.ID.String(), err)
			letter.MarkFailed()
			_ = s.referenceLetterRepo.Update(ctx, letter)
			continue
		}

		// Update profile summary if not set
		if profile.Summary == "" && profileData.Summary != "" {
			profile.Summary = profileData.Summary
		}

		// Create work experiences
		for _, expData := range profileData.WorkExperiences {
			experience := &domain.WorkExperience{
				ID:                uuid.New(),
				ProfileID:         profile.ID,
				CompanyName:       expData.CompanyName,
				Role:              expData.Role,
				Description:       expData.Description,
				ReferenceLetterID: &letter.ID,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}

			if expData.StartDate != nil {
				experience.StartDate = *expData.StartDate
			} else {
				experience.StartDate = time.Now() // Default to now if not provided
			}

			if expData.EndDate != nil {
				experience.EndDate = expData.EndDate
			}

			if err := s.workExperienceRepo.Create(ctx, experience); err != nil {
				s.logger.Error("Failed to create work experience: %v", err)
				continue
			}

			// Create credibility highlights for this experience
			for _, highlightData := range profileData.Highlights {
				highlight := domain.NewCredibilityHighlight(
					highlightData.Quote,
					domain.Sentiment(highlightData.Sentiment),
					letter.ID,
				)
				highlight.WorkExperienceID = experience.ID

				if err := s.highlightRepo.Create(ctx, highlight); err != nil {
					s.logger.Error("Failed to create credibility highlight: %v", err)
					continue
				}
			}
		}

		// Extract additional credibility highlights
		// Note: Currently these are tied to work experiences, so we skip profile-level highlights
		_, err = s.llmProvider.ExtractCredibilityHighlights(ctx, letter.ExtractedText)
		if err != nil {
			s.logger.Warn("Failed to extract additional credibility highlights: %v", err)
		}

		// Mark letter as processed
		letter.MarkProcessed(letter.ExtractedText)
		if err := s.referenceLetterRepo.Update(ctx, letter); err != nil {
			s.logger.Error("Failed to update reference letter status: %v", err)
		}
	}

	// Update profile
	profile.UpdatedAt = time.Now()
	if err := s.profileRepo.Update(ctx, profile); err != nil {
		s.logger.Error("Failed to update profile: %v", err)
		return nil, err
	}

	// Reload profile with relations
	profile, err = s.profileRepo.FindByUserIDWithRelations(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to reload profile: %v", err)
		return nil, err
	}

	s.logger.Info("Profile generation completed for user %s", userID.String())
	return profile, nil
}

// GetProfile retrieves a user's profile with all relations
func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	return s.profileRepo.FindByUserIDWithRelations(ctx, userID)
}

// UpdateProfile updates a user's profile
func (s *ProfileService) UpdateProfile(ctx context.Context, profile *domain.Profile) error {
	profile.UpdatedAt = time.Now()
	return s.profileRepo.Update(ctx, profile)
}
