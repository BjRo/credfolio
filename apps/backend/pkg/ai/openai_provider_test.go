package ai_test

import (
	"context"
	"testing"

	"github.com/credfolio/apps/backend/pkg/ai"
	"github.com/stretchr/testify/assert"
)

func TestNewOpenAIProvider_CreatesProviderWithAPIKey(t *testing.T) {
	// Act
	provider := ai.NewOpenAIProvider("test-api-key")

	// Assert
	assert.NotNil(t, provider)
}

// Note: Integration tests for OpenAI API calls are skipped in unit tests
// because they require real API calls. These would be tested in integration tests.

func TestOpenAIProvider_ExtractProfileData_WhenEmptyText_ReturnsError(t *testing.T) {
	// Arrange
	provider := ai.NewOpenAIProvider("test-api-key")

	// Act
	result, err := provider.ExtractProfileData(context.Background(), "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "reference text is required")
}

func TestOpenAIProvider_TailorProfile_WhenEmptyJobDescription_ReturnsError(t *testing.T) {
	// Arrange
	provider := ai.NewOpenAIProvider("test-api-key")

	// Act
	result, err := provider.TailorProfile(context.Background(), "summary", nil, nil, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "job description is required")
}

func TestOpenAIProvider_ExtractCredibilityHighlights_WhenEmptyText_ReturnsError(t *testing.T) {
	// Arrange
	provider := ai.NewOpenAIProvider("test-api-key")

	// Act
	result, err := provider.ExtractCredibilityHighlights(context.Background(), "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "reference text is required")
}
