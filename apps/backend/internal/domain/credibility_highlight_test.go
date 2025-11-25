package domain_test

import (
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredibilityHighlight_Validate_WhenValidData_ReturnsNoError(t *testing.T) {
	// Arrange
	sourceLetterID := uuid.New()
	ch := domain.NewCredibilityHighlight("Great team player", domain.SentimentPositive, sourceLetterID)

	// Act
	err := ch.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestCredibilityHighlight_Validate_WhenQuoteEmpty_ReturnsError(t *testing.T) {
	// Arrange
	sourceLetterID := uuid.New()
	ch := domain.NewCredibilityHighlight("", domain.SentimentPositive, sourceLetterID)

	// Act
	err := ch.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "quote is required")
}

func TestCredibilityHighlight_Validate_WhenSourceLetterIDEmpty_ReturnsError(t *testing.T) {
	// Arrange
	ch := domain.NewCredibilityHighlight("Great team player", domain.SentimentPositive, uuid.Nil)

	// Act
	err := ch.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source letter ID is required")
}

func TestCredibilityHighlight_Validate_WhenInvalidSentiment_ReturnsError(t *testing.T) {
	// Arrange
	sourceLetterID := uuid.New()
	ch := &domain.CredibilityHighlight{
		Quote:          "Great work",
		Sentiment:      "INVALID",
		SourceLetterID: sourceLetterID,
	}

	// Act
	err := ch.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "sentiment must be POSITIVE or NEUTRAL")
}

func TestCredibilityHighlight_IsPositive_WhenPositiveSentiment_ReturnsTrue(t *testing.T) {
	// Arrange
	sourceLetterID := uuid.New()
	ch := domain.NewCredibilityHighlight("Great team player", domain.SentimentPositive, sourceLetterID)

	// Act
	result := ch.IsPositive()

	// Assert
	assert.True(t, result)
}

func TestCredibilityHighlight_IsPositive_WhenNeutralSentiment_ReturnsFalse(t *testing.T) {
	// Arrange
	sourceLetterID := uuid.New()
	ch := domain.NewCredibilityHighlight("Worked on the project", domain.SentimentNeutral, sourceLetterID)

	// Act
	result := ch.IsPositive()

	// Assert
	assert.False(t, result)
}
