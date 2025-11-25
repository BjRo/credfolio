package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
