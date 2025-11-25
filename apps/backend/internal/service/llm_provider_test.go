package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/credfolio/apps/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

// MockLLMProvider is a mock implementation of LLMProvider for testing
type MockLLMProvider struct {
	ExtractProfileDataFunc           func(ctx context.Context, referenceText string) (*service.ProfileData, error)
	TailorProfileFunc                func(ctx context.Context, profileSummary string, experiences []service.WorkExperience, skills []string, jobDescription string) (*service.TailoredProfileData, error)
	ExtractCredibilityHighlightsFunc func(ctx context.Context, referenceText string) ([]service.HighlightData, error)
}

func (m *MockLLMProvider) ExtractProfileData(ctx context.Context, referenceText string) (*service.ProfileData, error) {
	if m.ExtractProfileDataFunc != nil {
		return m.ExtractProfileDataFunc(ctx, referenceText)
	}
	return nil, errors.New("not implemented")
}

func (m *MockLLMProvider) TailorProfile(ctx context.Context, profileSummary string, experiences []service.WorkExperience, skills []string, jobDescription string) (*service.TailoredProfileData, error) {
	if m.TailorProfileFunc != nil {
		return m.TailorProfileFunc(ctx, profileSummary, experiences, skills, jobDescription)
	}
	return nil, errors.New("not implemented")
}

func (m *MockLLMProvider) ExtractCredibilityHighlights(ctx context.Context, referenceText string) ([]service.HighlightData, error) {
	if m.ExtractCredibilityHighlightsFunc != nil {
		return m.ExtractCredibilityHighlightsFunc(ctx, referenceText)
	}
	return nil, errors.New("not implemented")
}

func TestMockLLMProvider_ExtractProfileData_WhenGivenValidText_ReturnsProfileData(t *testing.T) {
	// Arrange
	mockProvider := &MockLLMProvider{
		ExtractProfileDataFunc: func(_ context.Context, referenceText string) (*service.ProfileData, error) {
			return &service.ProfileData{
				Summary: "Experienced software engineer",
				WorkExperiences: []service.WorkExperience{
					{CompanyName: "Test Corp", Role: "Senior Engineer"},
				},
				Skills: []string{"Go", "Python"},
			}, nil
		},
	}

	// Act
	result, err := mockProvider.ExtractProfileData(context.Background(), "Sample reference letter text")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Experienced software engineer", result.Summary)
	assert.Len(t, result.WorkExperiences, 1)
	assert.Len(t, result.Skills, 2)
}

func TestMockLLMProvider_ExtractProfileData_WhenProviderFails_ReturnsError(t *testing.T) {
	// Arrange
	mockProvider := &MockLLMProvider{
		ExtractProfileDataFunc: func(_ context.Context, _ string) (*service.ProfileData, error) {
			return nil, errors.New("API error")
		},
	}

	// Act
	result, err := mockProvider.ExtractProfileData(context.Background(), "Sample text")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "API error")
}

func TestMockLLMProvider_TailorProfile_WhenGivenValidInput_ReturnsTailoredData(t *testing.T) {
	// Arrange
	mockProvider := &MockLLMProvider{
		TailorProfileFunc: func(_ context.Context, _ string, _ []service.WorkExperience, _ []string, _ string) (*service.TailoredProfileData, error) {
			return &service.TailoredProfileData{
				TailoredSummary:   "Tailored for Go development position",
				MatchScore:        0.85,
				HighlightedSkills: []string{"Go", "Kubernetes"},
			}, nil
		},
	}

	// Act
	result, err := mockProvider.TailorProfile(
		context.Background(),
		"Original summary",
		[]service.WorkExperience{},
		[]string{"Go", "Python"},
		"Looking for Go developer",
	)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0.85, result.MatchScore)
}

func TestMockLLMProvider_ExtractCredibilityHighlights_WhenGivenText_ReturnsHighlights(t *testing.T) {
	// Arrange
	mockProvider := &MockLLMProvider{
		ExtractCredibilityHighlightsFunc: func(_ context.Context, _ string) ([]service.HighlightData, error) {
			return []service.HighlightData{
				{Quote: "Excellent team player", Sentiment: "POSITIVE"},
				{Quote: "Strong technical skills", Sentiment: "POSITIVE"},
			}, nil
		},
	}

	// Act
	result, err := mockProvider.ExtractCredibilityHighlights(context.Background(), "Reference text")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "POSITIVE", result[0].Sentiment)
}
