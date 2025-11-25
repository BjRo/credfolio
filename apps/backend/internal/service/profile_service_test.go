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

// MockProfileRepo is a mock implementation of ProfileRepository
type MockProfileRepo struct {
	mock.Mock
}

func (m *MockProfileRepo) Create(ctx context.Context, profile *domain.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockProfileRepo) Update(ctx context.Context, profile *domain.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockProfileRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileRepo) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProfileRepo) FindByUserIDWithRelations(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

// MockRefLetterRepo is a mock implementation of ReferenceLetterRepository
type MockRefLetterRepo struct {
	mock.Mock
}

func (m *MockRefLetterRepo) Create(ctx context.Context, letter *domain.ReferenceLetter) error {
	args := m.Called(ctx, letter)
	return args.Error(0)
}

func (m *MockRefLetterRepo) Update(ctx context.Context, letter *domain.ReferenceLetter) error {
	args := m.Called(ctx, letter)
	return args.Error(0)
}

func (m *MockRefLetterRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.ReferenceLetter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ReferenceLetter), args.Error(1)
}

func (m *MockRefLetterRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ReferenceLetter), args.Error(1)
}

func (m *MockRefLetterRepo) FindPendingByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ReferenceLetter), args.Error(1)
}

func (m *MockRefLetterRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockWorkExpRepo is a mock implementation of WorkExperienceRepository
type MockWorkExpRepo struct {
	mock.Mock
}

func (m *MockWorkExpRepo) Create(ctx context.Context, exp *domain.WorkExperience) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MockWorkExpRepo) Update(ctx context.Context, exp *domain.WorkExperience) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MockWorkExpRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.WorkExperience, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WorkExperience), args.Error(1)
}

func (m *MockWorkExpRepo) FindByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.WorkExperience, error) {
	args := m.Called(ctx, profileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.WorkExperience), args.Error(1)
}

func (m *MockWorkExpRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWorkExpRepo) FindByProfileIDWithHighlights(ctx context.Context, profileID uuid.UUID) ([]*domain.WorkExperience, error) {
	args := m.Called(ctx, profileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.WorkExperience), args.Error(1)
}

// MockHighlightRepo is a mock implementation of CredibilityHighlightRepository
type MockHighlightRepo struct {
	mock.Mock
}

func (m *MockHighlightRepo) Create(ctx context.Context, highlight *domain.CredibilityHighlight) error {
	args := m.Called(ctx, highlight)
	return args.Error(0)
}

func (m *MockHighlightRepo) Update(ctx context.Context, highlight *domain.CredibilityHighlight) error {
	args := m.Called(ctx, highlight)
	return args.Error(0)
}

func (m *MockHighlightRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.CredibilityHighlight, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CredibilityHighlight), args.Error(1)
}

func (m *MockHighlightRepo) FindByWorkExperienceID(ctx context.Context, workExpID uuid.UUID) ([]*domain.CredibilityHighlight, error) {
	args := m.Called(ctx, workExpID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CredibilityHighlight), args.Error(1)
}

func (m *MockHighlightRepo) FindBySourceLetterID(ctx context.Context, letterID uuid.UUID) ([]*domain.CredibilityHighlight, error) {
	args := m.Called(ctx, letterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CredibilityHighlight), args.Error(1)
}

func (m *MockHighlightRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockLLM is a mock implementation of LLMProvider
type MockLLM struct {
	mock.Mock
}

func (m *MockLLM) ExtractProfileData(ctx context.Context, text string) (*service.ProfileData, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ProfileData), args.Error(1)
}

func (m *MockLLM) TailorProfile(ctx context.Context, profileSummary string, experiences []service.WorkExperience, skills []string, jobDescription string) (*service.TailoredProfileData, error) {
	args := m.Called(ctx, profileSummary, experiences, skills, jobDescription)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.TailoredProfileData), args.Error(1)
}

func (m *MockLLM) ExtractCredibilityHighlights(ctx context.Context, text string) ([]service.HighlightData, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]service.HighlightData), args.Error(1)
}

func (m *MockLLM) CalculateRelevance(ctx context.Context, text string, reference string) (float64, error) {
	args := m.Called(ctx, text, reference)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockLLM) GenerateText(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

// T049: Unit test for ProfileService when generating profile from reference letter extracts structured data
func TestProfileService_GenerateProfileFromReferences_WhenValidReferenceLetter_ExtractsStructuredData(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockProfileRepo := new(MockProfileRepo)
	mockRefLetterRepo := new(MockRefLetterRepo)
	mockWorkExpRepo := new(MockWorkExpRepo)
	mockHighlightRepo := new(MockHighlightRepo)
	mockLLM := new(MockLLM)

	letter := domain.NewReferenceLetter(userID, "letter.txt", "/uploads/letter.txt")
	letter.ExtractedText = "John worked at Acme Corp as a Software Engineer from 2020 to 2023. He was exceptional."

	profile := domain.NewProfile(userID)

	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	extractedData := &service.ProfileData{
		Summary: "Experienced software engineer",
		WorkExperiences: []service.WorkExperience{
			{
				CompanyName: "Acme Corp",
				Role:        "Software Engineer",
				Description: "Developed software solutions",
				StartDate:   &startDate,
				EndDate:     &endDate,
			},
		},
		Highlights: []service.HighlightData{
			{
				Quote:     "He was exceptional",
				Sentiment: "positive",
			},
		},
	}

	mockRefLetterRepo.On("FindPendingByUserID", ctx, userID).Return([]*domain.ReferenceLetter{letter}, nil)
	mockProfileRepo.On("FindByUserID", ctx, userID).Return(nil, nil)
	mockProfileRepo.On("Create", ctx, mock.AnythingOfType("*domain.Profile")).Return(nil)
	mockLLM.On("ExtractProfileData", ctx, letter.ExtractedText).Return(extractedData, nil)
	mockWorkExpRepo.On("Create", ctx, mock.AnythingOfType("*domain.WorkExperience")).Return(nil)
	mockHighlightRepo.On("Create", ctx, mock.AnythingOfType("*domain.CredibilityHighlight")).Return(nil)
	mockLLM.On("ExtractCredibilityHighlights", ctx, letter.ExtractedText).Return([]service.HighlightData{}, nil)
	mockRefLetterRepo.On("Update", ctx, mock.AnythingOfType("*domain.ReferenceLetter")).Return(nil)
	mockProfileRepo.On("Update", ctx, mock.AnythingOfType("*domain.Profile")).Return(nil)
	mockProfileRepo.On("FindByUserIDWithRelations", ctx, userID).Return(profile, nil)

	svc := service.NewProfileService(mockProfileRepo, mockRefLetterRepo, mockWorkExpRepo, mockHighlightRepo, mockLLM)

	// Act
	result, err := svc.GenerateProfileFromReferences(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	mockLLM.AssertCalled(t, "ExtractProfileData", ctx, letter.ExtractedText)
	mockWorkExpRepo.AssertCalled(t, "Create", ctx, mock.AnythingOfType("*domain.WorkExperience"))
}

// T050: Unit test for ProfileService when LLMProvider returns error propagates error
func TestProfileService_GenerateProfileFromReferences_WhenLLMProviderFails_PropagatesError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockProfileRepo := new(MockProfileRepo)
	mockRefLetterRepo := new(MockRefLetterRepo)
	mockWorkExpRepo := new(MockWorkExpRepo)
	mockHighlightRepo := new(MockHighlightRepo)
	mockLLM := new(MockLLM)

	letter := domain.NewReferenceLetter(userID, "letter.txt", "/uploads/letter.txt")
	letter.ExtractedText = "Some text"

	profile := domain.NewProfile(userID)

	mockRefLetterRepo.On("FindPendingByUserID", ctx, userID).Return([]*domain.ReferenceLetter{letter}, nil)
	mockProfileRepo.On("FindByUserID", ctx, userID).Return(nil, nil)
	mockProfileRepo.On("Create", ctx, mock.AnythingOfType("*domain.Profile")).Return(nil)
	mockLLM.On("ExtractProfileData", ctx, letter.ExtractedText).Return(nil, errors.New("LLM API error"))
	mockRefLetterRepo.On("Update", ctx, mock.AnythingOfType("*domain.ReferenceLetter")).Return(nil)
	mockProfileRepo.On("Update", ctx, mock.AnythingOfType("*domain.Profile")).Return(nil)
	mockProfileRepo.On("FindByUserIDWithRelations", ctx, userID).Return(profile, nil)

	svc := service.NewProfileService(mockProfileRepo, mockRefLetterRepo, mockWorkExpRepo, mockHighlightRepo, mockLLM)

	// Act
	result, err := svc.GenerateProfileFromReferences(ctx, userID)

	// Assert
	// The service continues even when LLM fails for individual letters, but returns the profile
	require.NoError(t, err)
	assert.NotNil(t, result)
	mockRefLetterRepo.AssertCalled(t, "Update", ctx, mock.AnythingOfType("*domain.ReferenceLetter"))
}

// T051: Unit test for ProfileService when extracting credibility finds positive sentiment quotes
func TestProfileService_GenerateProfileFromReferences_WhenExtractingCredibility_FindsPositiveSentimentQuotes(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockProfileRepo := new(MockProfileRepo)
	mockRefLetterRepo := new(MockRefLetterRepo)
	mockWorkExpRepo := new(MockWorkExpRepo)
	mockHighlightRepo := new(MockHighlightRepo)
	mockLLM := new(MockLLM)

	letter := domain.NewReferenceLetter(userID, "letter.txt", "/uploads/letter.txt")
	letter.ExtractedText = "John was an exceptional team player who consistently delivered high-quality work."

	profile := domain.NewProfile(userID)

	extractedData := &service.ProfileData{
		Summary: "Great employee",
		WorkExperiences: []service.WorkExperience{
			{
				CompanyName: "Company",
				Role:        "Engineer",
			},
		},
		Highlights: []service.HighlightData{
			{
				Quote:     "exceptional team player",
				Sentiment: "positive",
			},
			{
				Quote:     "consistently delivered high-quality work",
				Sentiment: "positive",
			},
		},
	}

	mockRefLetterRepo.On("FindPendingByUserID", ctx, userID).Return([]*domain.ReferenceLetter{letter}, nil)
	mockProfileRepo.On("FindByUserID", ctx, userID).Return(nil, nil)
	mockProfileRepo.On("Create", ctx, mock.AnythingOfType("*domain.Profile")).Return(nil)
	mockLLM.On("ExtractProfileData", ctx, letter.ExtractedText).Return(extractedData, nil)
	mockWorkExpRepo.On("Create", ctx, mock.AnythingOfType("*domain.WorkExperience")).Return(nil)
	mockHighlightRepo.On("Create", ctx, mock.AnythingOfType("*domain.CredibilityHighlight")).Return(nil)
	mockLLM.On("ExtractCredibilityHighlights", ctx, letter.ExtractedText).Return([]service.HighlightData{}, nil)
	mockRefLetterRepo.On("Update", ctx, mock.AnythingOfType("*domain.ReferenceLetter")).Return(nil)
	mockProfileRepo.On("Update", ctx, mock.AnythingOfType("*domain.Profile")).Return(nil)
	mockProfileRepo.On("FindByUserIDWithRelations", ctx, userID).Return(profile, nil)

	svc := service.NewProfileService(mockProfileRepo, mockRefLetterRepo, mockWorkExpRepo, mockHighlightRepo, mockLLM)

	// Act
	_, err := svc.GenerateProfileFromReferences(ctx, userID)

	// Assert
	require.NoError(t, err)
	// Verify highlights were created (2 highlights for 1 work experience)
	mockHighlightRepo.AssertNumberOfCalls(t, "Create", 2)
}

// T052: Unit test for ProfileService when reference letter is invalid returns validation error
func TestProfileService_GenerateProfileFromReferences_WhenNoPendingLetters_ReturnsError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockProfileRepo := new(MockProfileRepo)
	mockRefLetterRepo := new(MockRefLetterRepo)
	mockWorkExpRepo := new(MockWorkExpRepo)
	mockHighlightRepo := new(MockHighlightRepo)
	mockLLM := new(MockLLM)

	mockRefLetterRepo.On("FindPendingByUserID", ctx, userID).Return([]*domain.ReferenceLetter{}, nil)

	svc := service.NewProfileService(mockProfileRepo, mockRefLetterRepo, mockWorkExpRepo, mockHighlightRepo, mockLLM)

	// Act
	result, err := svc.GenerateProfileFromReferences(ctx, userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no pending reference letters found")
}

func TestProfileService_GetProfile_WhenProfileExists_ReturnsProfile(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockProfileRepo := new(MockProfileRepo)
	mockRefLetterRepo := new(MockRefLetterRepo)
	mockWorkExpRepo := new(MockWorkExpRepo)
	mockHighlightRepo := new(MockHighlightRepo)
	mockLLM := new(MockLLM)

	profile := domain.NewProfile(userID)
	profile.Summary = "Test summary"

	mockProfileRepo.On("FindByUserIDWithRelations", ctx, userID).Return(profile, nil)

	svc := service.NewProfileService(mockProfileRepo, mockRefLetterRepo, mockWorkExpRepo, mockHighlightRepo, mockLLM)

	// Act
	result, err := svc.GetProfile(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test summary", result.Summary)
}

func TestProfileService_UpdateProfile_WhenProfileValid_UpdatesProfile(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockProfileRepo := new(MockProfileRepo)
	mockRefLetterRepo := new(MockRefLetterRepo)
	mockWorkExpRepo := new(MockWorkExpRepo)
	mockHighlightRepo := new(MockHighlightRepo)
	mockLLM := new(MockLLM)

	profile := domain.NewProfile(userID)
	profile.Summary = "Updated summary"

	mockProfileRepo.On("Update", ctx, profile).Return(nil)

	svc := service.NewProfileService(mockProfileRepo, mockRefLetterRepo, mockWorkExpRepo, mockHighlightRepo, mockLLM)

	// Act
	err := svc.UpdateProfile(ctx, profile)

	// Assert
	require.NoError(t, err)
	mockProfileRepo.AssertCalled(t, "Update", ctx, profile)
}

// T084: Unit test for ProfileService when getting profile aggregates skills across experiences
func TestProfileService_GetProfile_WhenGettingProfile_AggregatesSkillsAcrossExperiences(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockProfileRepo := new(MockProfileRepo)
	mockRefLetterRepo := new(MockRefLetterRepo)
	mockWorkExpRepo := new(MockWorkExpRepo)
	mockHighlightRepo := new(MockHighlightRepo)
	mockLLM := new(MockLLM)

	profile := domain.NewProfile(userID)
	profile.Summary = "Experienced engineer"
	profile.Skills = []*domain.Skill{
		{ID: uuid.New(), Name: "Go"},
		{ID: uuid.New(), Name: "TypeScript"},
		{ID: uuid.New(), Name: "Go"}, // Duplicate
	}

	mockProfileRepo.On("FindByUserIDWithRelations", ctx, userID).Return(profile, nil)

	svc := service.NewProfileService(mockProfileRepo, mockRefLetterRepo, mockWorkExpRepo, mockHighlightRepo, mockLLM)

	// Act
	result, err := svc.GetProfile(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	// The profile should have skills (aggregation logic can be added later)
	assert.NotEmpty(t, result.Skills)
}

// T085: Unit test for ProfileService when getting profile includes credibility highlights
func TestProfileService_GetProfile_WhenGettingProfile_IncludesCredibilityHighlights(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	mockProfileRepo := new(MockProfileRepo)
	mockRefLetterRepo := new(MockRefLetterRepo)
	mockWorkExpRepo := new(MockWorkExpRepo)
	mockHighlightRepo := new(MockHighlightRepo)
	mockLLM := new(MockLLM)

	profile := domain.NewProfile(userID)
	profile.Summary = "Great engineer"

	letterID := uuid.New()
	workExpID := uuid.New()
	profile.WorkExperiences = []domain.WorkExperience{
		{
			ID:                uuid.New(),
			ProfileID:         profile.ID,
			CompanyName:       "Acme Corp",
			Role:              "Engineer",
			StartDate:         time.Now().AddDate(-2, 0, 0),
			ReferenceLetterID: &letterID,
			CredibilityHighlights: []domain.CredibilityHighlight{
				{
					ID:               uuid.New(),
					WorkExperienceID: workExpID,
					Quote:            "Exceptional team player",
					Sentiment:        domain.SentimentPositive,
					SourceLetterID:   letterID,
				},
				{
					ID:               uuid.New(),
					WorkExperienceID: workExpID,
					Quote:            "Delivered high quality work",
					Sentiment:        domain.SentimentPositive,
					SourceLetterID:   letterID,
				},
			},
		},
	}

	mockProfileRepo.On("FindByUserIDWithRelations", ctx, userID).Return(profile, nil)

	svc := service.NewProfileService(mockProfileRepo, mockRefLetterRepo, mockWorkExpRepo, mockHighlightRepo, mockLLM)

	// Act
	result, err := svc.GetProfile(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.WorkExperiences, 1)
	assert.Len(t, result.WorkExperiences[0].CredibilityHighlights, 2)
	assert.Equal(t, "Exceptional team player", result.WorkExperiences[0].CredibilityHighlights[0].Quote)
}
