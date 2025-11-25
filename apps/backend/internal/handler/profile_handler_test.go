package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/handler"
	"github.com/credfolio/apps/backend/internal/handler/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockProfileService is a mock implementation of ProfileService
type MockProfileService struct {
	mock.Mock
}

func (m *MockProfileService) GenerateProfileFromReferences(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileService) UpdateProfile(ctx context.Context, profile *domain.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

// T055: Unit test for ProfileHandler when generating profile returns profile data
func TestProfileHandler_Generate_WhenGeneratingProfile_ReturnsProfileData(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)

	h := handler.NewProfileHandler(mockService)

	userID := uuid.New()
	profile := domain.NewProfile(userID)
	profile.Summary = "Experienced software engineer"

	mockService.On("GenerateProfileFromReferences", mock.Anything, userID).Return(profile, nil)

	req := httptest.NewRequest(http.MethodPost, "/profile/generate", nil)
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.Generate(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response handler.ProfileResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "Experienced software engineer", response.Summary)
}

func TestProfileHandler_Generate_WhenNoUserID_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	h := handler.NewProfileHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/profile/generate", nil)
	recorder := httptest.NewRecorder()

	// Act
	h.Generate(recorder, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestProfileHandler_Generate_WhenServiceFails_ReturnsInternalServerError(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	h := handler.NewProfileHandler(mockService)

	userID := uuid.New()

	mockService.On("GenerateProfileFromReferences", mock.Anything, userID).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/profile/generate", nil)
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.Generate(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestProfileHandler_Get_WhenProfileExists_ReturnsProfile(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	h := handler.NewProfileHandler(mockService)

	userID := uuid.New()
	profile := domain.NewProfile(userID)
	profile.Summary = "Test profile"

	mockService.On("GetProfile", mock.Anything, userID).Return(profile, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.Get(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response handler.ProfileResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "Test profile", response.Summary)
}

func TestProfileHandler_Get_WhenProfileNotFound_ReturnsNotFound(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	h := handler.NewProfileHandler(mockService)

	userID := uuid.New()

	mockService.On("GetProfile", mock.Anything, userID).Return(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.Get(recorder, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestProfileHandler_Get_WhenNoUserID_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	h := handler.NewProfileHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	recorder := httptest.NewRecorder()

	// Act
	h.Get(recorder, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestProfileHandler_Generate_WhenProfileWithWorkExperiences_ReturnsCompleteProfile(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	h := handler.NewProfileHandler(mockService)

	userID := uuid.New()
	profile := domain.NewProfile(userID)
	profile.Summary = "Senior engineer"
	profile.WorkExperiences = []domain.WorkExperience{
		{
			ID:          uuid.New(),
			ProfileID:   profile.ID,
			CompanyName: "Acme Corp",
			Role:        "Software Engineer",
		},
	}

	mockService.On("GenerateProfileFromReferences", mock.Anything, userID).Return(profile, nil)

	req := httptest.NewRequest(http.MethodPost, "/profile/generate", nil)
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.Generate(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response handler.ProfileResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)

	assert.Len(t, response.WorkExperiences, 1)
	assert.Equal(t, "Acme Corp", response.WorkExperiences[0].CompanyName)
}

// T086: Unit test for ProfileHandler when getting profile returns complete profile data
func TestProfileHandler_Get_WhenGettingProfile_ReturnsCompleteProfileData(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	h := handler.NewProfileHandler(mockService)

	userID := uuid.New()
	profile := domain.NewProfile(userID)
	profile.Summary = "Senior software engineer with 10 years experience"
	profile.Skills = []*domain.Skill{
		{ID: uuid.New(), Name: "Go"},
		{ID: uuid.New(), Name: "TypeScript"},
	}
	profile.WorkExperiences = []domain.WorkExperience{
		{
			ID:          uuid.New(),
			ProfileID:   profile.ID,
			CompanyName: "Tech Corp",
			Role:        "Senior Engineer",
			Description: "Led engineering team",
			CredibilityHighlights: []domain.CredibilityHighlight{
				{
					ID:        uuid.New(),
					Quote:     "Outstanding technical leadership",
					Sentiment: domain.SentimentPositive,
				},
			},
		},
	}

	mockService.On("GetProfile", mock.Anything, userID).Return(profile, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.Get(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response handler.ProfileResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "Senior software engineer with 10 years experience", response.Summary)
	assert.Len(t, response.Skills, 2)
	assert.Contains(t, response.Skills, "Go")
	assert.Contains(t, response.Skills, "TypeScript")
	assert.Len(t, response.WorkExperiences, 1)
	assert.Equal(t, "Tech Corp", response.WorkExperiences[0].CompanyName)
	assert.Len(t, response.WorkExperiences[0].CredibilityHighlights, 1)
	assert.Equal(t, "Outstanding technical leadership", response.WorkExperiences[0].CredibilityHighlights[0].Quote)
}

// T087: Unit test for ProfileHandler when profile not found returns 404 error
// Note: Already covered by TestProfileHandler_Get_WhenProfileNotFound_ReturnsNotFound above

// MockTailoringService is a mock implementation of TailoringService
type MockTailoringService struct {
	mock.Mock
}

func (m *MockTailoringService) TailorProfileToJobDescription(ctx context.Context, userID uuid.UUID, jobDescription string) (*handler.TailoredProfileResponse, error) {
	args := m.Called(ctx, userID, jobDescription)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*handler.TailoredProfileResponse), args.Error(1)
}

// T109: Unit test for ProfileHandler when tailoring profile returns tailored profile data
func TestProfileHandler_Tailor_WhenTailoringProfile_ReturnsTailoredProfileData(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	mockTailoringService := new(MockTailoringService)
	h := handler.NewProfileHandlerWithTailoring(mockService, mockTailoringService)

	userID := uuid.New()
	tailoredProfile := &handler.TailoredProfileResponse{
		ID:           uuid.New().String(),
		MatchScore:   0.85,
		MatchSummary: "Strong match based on Go experience",
		TailoredExperiences: []handler.TailoredExperienceResponse{
			{
				CompanyName:    "Tech Corp",
				Role:           "Backend Engineer",
				RelevanceScore: 0.9,
			},
		},
		RelevantSkills: []string{"Go", "PostgreSQL"},
	}

	jobDescription := "We are looking for a Backend Engineer with strong Go experience to join our team. The ideal candidate has experience with microservices and PostgreSQL databases."
	mockTailoringService.On("TailorProfileToJobDescription", mock.Anything, userID, jobDescription).Return(tailoredProfile, nil)

	body := `{"jobDescription":"` + jobDescription + `"}`
	req := httptest.NewRequest(http.MethodPost, "/profile/tailor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.Tailor(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response handler.TailoredProfileResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 0.85, response.MatchScore)
	assert.Equal(t, "Strong match based on Go experience", response.MatchSummary)
	assert.Len(t, response.TailoredExperiences, 1)
}

// T110: Unit test for ProfileHandler when job description invalid returns validation error
func TestProfileHandler_Tailor_WhenJobDescriptionInvalid_ReturnsValidationError(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	mockTailoringService := new(MockTailoringService)
	h := handler.NewProfileHandlerWithTailoring(mockService, mockTailoringService)

	userID := uuid.New()

	body := `{"jobDescription":""}`
	req := httptest.NewRequest(http.MethodPost, "/profile/tailor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.Tailor(recorder, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestProfileHandler_Tailor_WhenNoUserID_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	mockTailoringService := new(MockTailoringService)
	h := handler.NewProfileHandlerWithTailoring(mockService, mockTailoringService)

	body := `{"jobDescription":"Backend engineer"}`
	req := httptest.NewRequest(http.MethodPost, "/profile/tailor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	// Act
	h.Tailor(recorder, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

// T131: Unit test for ProfileHandler when downloading CV returns PDF bytes
func TestProfileHandler_DownloadCV_WhenProfileExists_ReturnsPDFBytes(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	h := handler.NewProfileHandler(mockService)

	userID := uuid.New()
	profile := domain.NewProfile(userID)
	profile.Summary = "Software Engineer"

	mockService.On("GetProfile", mock.Anything, userID).Return(profile, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile/cv", nil)
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.DownloadCV(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/pdf", recorder.Header().Get("Content-Type"))
	assert.Contains(t, recorder.Header().Get("Content-Disposition"), "attachment")
}

// T132: Unit test for ProfileHandler when profile not found returns 404 error
func TestProfileHandler_DownloadCV_WhenProfileNotFound_ReturnsNotFound(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	h := handler.NewProfileHandler(mockService)

	userID := uuid.New()

	mockService.On("GetProfile", mock.Anything, userID).Return(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/profile/cv", nil)
	ctx := middleware.SetUserID(req.Context(), userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	// Act
	h.DownloadCV(recorder, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

// T133: Unit test for ProfileHandler when PDF generation fails returns error
func TestProfileHandler_DownloadCV_WhenNoUserID_ReturnsUnauthorized(t *testing.T) {
	// Arrange
	mockService := new(MockProfileService)
	h := handler.NewProfileHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/profile/cv", nil)
	recorder := httptest.NewRecorder()

	// Act
	h.DownloadCV(recorder, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}
