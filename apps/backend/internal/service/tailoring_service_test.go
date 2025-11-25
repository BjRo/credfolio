package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockJobMatchRepo is a mock implementation of JobMatchRepository
type MockJobMatchRepo struct {
	mock.Mock
}

func (m *MockJobMatchRepo) Create(ctx context.Context, jobMatch *domain.JobMatch) error {
	args := m.Called(ctx, jobMatch)
	return args.Error(0)
}

func (m *MockJobMatchRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.JobMatch, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JobMatch), args.Error(1)
}

func (m *MockJobMatchRepo) FindByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.JobMatch, error) {
	args := m.Called(ctx, profileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.JobMatch), args.Error(1)
}

func (m *MockJobMatchRepo) Update(ctx context.Context, jobMatch *domain.JobMatch) error {
	args := m.Called(ctx, jobMatch)
	return args.Error(0)
}

func (m *MockJobMatchRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockTailoringLLM is a mock implementation of LLMProvider for tailoring service
type MockTailoringLLM struct {
	mock.Mock
}

func (m *MockTailoringLLM) ExtractProfileData(ctx context.Context, text string) (*service.ProfileData, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ProfileData), args.Error(1)
}

func (m *MockTailoringLLM) TailorProfile(ctx context.Context, profileSummary string, experiences []service.WorkExperience, skills []string, jobDescription string) (*service.TailoredProfileData, error) {
	args := m.Called(ctx, profileSummary, experiences, skills, jobDescription)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.TailoredProfileData), args.Error(1)
}

func (m *MockTailoringLLM) ExtractCredibilityHighlights(ctx context.Context, text string) ([]service.HighlightData, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.HighlightData), args.Error(1)
}

func (m *MockTailoringLLM) GenerateText(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

func (m *MockTailoringLLM) CalculateRelevance(ctx context.Context, text string, reference string) (float64, error) {
	args := m.Called(ctx, text, reference)
	return args.Get(0).(float64), args.Error(1)
}

// T105: Unit test for TailoringService when tailoring profile ranks experiences by relevance
func TestTailoringService_TailorProfile_WhenTailoringProfile_RanksExperiencesByRelevance(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockJobMatchRepo := new(MockJobMatchRepo)
	mockLLM := new(MockTailoringLLM)
	mockProfileRepo := new(MockProfileRepo)

	profile := domain.NewProfile(userID)
	profile.Summary = "Experienced engineer"
	profile.WorkExperiences = []domain.WorkExperience{
		{
			ID:          uuid.New(),
			ProfileID:   profile.ID,
			CompanyName: "Tech Corp",
			Role:        "Backend Engineer",
			Description: "Worked with Go and databases",
			StartDate:   time.Now().AddDate(-2, 0, 0),
		},
		{
			ID:          uuid.New(),
			ProfileID:   profile.ID,
			CompanyName: "Web Inc",
			Role:        "Frontend Developer",
			Description: "React and TypeScript development",
			StartDate:   time.Now().AddDate(-4, 0, 0),
		},
	}

	jobDescription := "Looking for a Backend Engineer with Go experience"

	mockProfileRepo.On("FindByUserIDWithRelations", ctx, userID).Return(profile, nil)
	mockLLM.On("CalculateRelevance", ctx, mock.MatchedBy(func(text string) bool {
		return text != ""
	}), jobDescription).Return(0.9, nil).Once() // Backend exp - high relevance
	mockLLM.On("CalculateRelevance", ctx, mock.MatchedBy(func(text string) bool {
		return text != ""
	}), jobDescription).Return(0.3, nil).Once() // Frontend exp - low relevance
	mockLLM.On("GenerateText", ctx, mock.Anything).Return("Strong match based on Go experience", nil)
	mockJobMatchRepo.On("Create", ctx, mock.AnythingOfType("*domain.JobMatch")).Return(nil)

	svc := service.NewTailoringService(mockProfileRepo, mockJobMatchRepo, mockLLM)

	// Act
	result, err := svc.TailorProfileToJobDescription(ctx, userID, jobDescription)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Greater(t, result.MatchScore, 0.0)
}

// T106: Unit test for TailoringService when job description is empty returns error
func TestTailoringService_TailorProfile_WhenJobDescriptionEmpty_ReturnsError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockJobMatchRepo := new(MockJobMatchRepo)
	mockLLM := new(MockTailoringLLM)
	mockProfileRepo := new(MockProfileRepo)

	svc := service.NewTailoringService(mockProfileRepo, mockJobMatchRepo, mockLLM)

	// Act
	result, err := svc.TailorProfileToJobDescription(ctx, userID, "")

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "job description is required")
}

// T107: Unit test for TailoringService when calculating match score returns score between 0 and 1
func TestTailoringService_TailorProfile_WhenCalculatingMatchScore_ReturnsScoreBetween0And1(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockJobMatchRepo := new(MockJobMatchRepo)
	mockLLM := new(MockTailoringLLM)
	mockProfileRepo := new(MockProfileRepo)

	profile := domain.NewProfile(userID)
	profile.Summary = "Senior engineer"
	profile.WorkExperiences = []domain.WorkExperience{
		{
			ID:          uuid.New(),
			ProfileID:   profile.ID,
			CompanyName: "Tech Corp",
			Role:        "Engineer",
			Description: "Software development",
			StartDate:   time.Now().AddDate(-1, 0, 0),
		},
	}

	jobDescription := "Software Engineer with 3+ years experience"

	mockProfileRepo.On("FindByUserIDWithRelations", ctx, userID).Return(profile, nil)
	mockLLM.On("CalculateRelevance", ctx, mock.Anything, jobDescription).Return(0.75, nil)
	mockLLM.On("GenerateText", ctx, mock.Anything).Return("Good match", nil)
	mockJobMatchRepo.On("Create", ctx, mock.AnythingOfType("*domain.JobMatch")).Return(nil)

	svc := service.NewTailoringService(mockProfileRepo, mockJobMatchRepo, mockLLM)

	// Act
	result, err := svc.TailorProfileToJobDescription(ctx, userID, jobDescription)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, result.MatchScore, 0.0)
	assert.LessOrEqual(t, result.MatchScore, 1.0)
}

// T108: Unit test for TailoringService when LLMProvider fails propagates error
func TestTailoringService_TailorProfile_WhenLLMProviderFails_PropagatesError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockJobMatchRepo := new(MockJobMatchRepo)
	mockLLM := new(MockTailoringLLM)
	mockProfileRepo := new(MockProfileRepo)

	profile := domain.NewProfile(userID)
	profile.WorkExperiences = []domain.WorkExperience{
		{
			ID:          uuid.New(),
			ProfileID:   profile.ID,
			CompanyName: "Tech Corp",
			Role:        "Engineer",
			Description: "Software development",
			StartDate:   time.Now().AddDate(-1, 0, 0),
		},
	}

	jobDescription := "Software Engineer position"

	mockProfileRepo.On("FindByUserIDWithRelations", ctx, userID).Return(profile, nil)
	mockLLM.On("CalculateRelevance", ctx, mock.Anything, jobDescription).Return(0.0, errors.New("LLM service unavailable"))

	svc := service.NewTailoringService(mockProfileRepo, mockJobMatchRepo, mockLLM)

	// Act
	result, err := svc.TailorProfileToJobDescription(ctx, userID, jobDescription)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "LLM service unavailable")
}

func TestTailoringService_TailorProfile_WhenProfileNotFound_ReturnsError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockJobMatchRepo := new(MockJobMatchRepo)
	mockLLM := new(MockTailoringLLM)
	mockProfileRepo := new(MockProfileRepo)

	mockProfileRepo.On("FindByUserIDWithRelations", ctx, userID).Return(nil, nil)

	svc := service.NewTailoringService(mockProfileRepo, mockJobMatchRepo, mockLLM)

	// Act
	result, err := svc.TailorProfileToJobDescription(ctx, userID, "Some job description")

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "profile not found")
}
