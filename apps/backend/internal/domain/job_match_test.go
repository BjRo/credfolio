package domain_test

import (
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobMatch_Validate_WhenValidData_ReturnsNoError(t *testing.T) {
	// Arrange
	profileID := uuid.New()
	jm := domain.NewJobMatch(profileID, "Looking for a Go developer...")
	jm.MatchScore = 0.85

	// Act
	err := jm.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestJobMatch_Validate_WhenBaseProfileIDEmpty_ReturnsError(t *testing.T) {
	// Arrange
	jm := domain.NewJobMatch(uuid.Nil, "Looking for a Go developer...")

	// Act
	err := jm.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "base profile ID is required")
}

func TestJobMatch_Validate_WhenJobDescriptionEmpty_ReturnsError(t *testing.T) {
	// Arrange
	profileID := uuid.New()
	jm := domain.NewJobMatch(profileID, "")

	// Act
	err := jm.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "job description is required")
}

func TestJobMatch_Validate_WhenMatchScoreNegative_ReturnsError(t *testing.T) {
	// Arrange
	profileID := uuid.New()
	jm := domain.NewJobMatch(profileID, "Looking for a Go developer...")
	jm.MatchScore = -0.1

	// Act
	err := jm.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "match score must be between 0 and 1")
}

func TestJobMatch_Validate_WhenMatchScoreAboveOne_ReturnsError(t *testing.T) {
	// Arrange
	profileID := uuid.New()
	jm := domain.NewJobMatch(profileID, "Looking for a Go developer...")
	jm.MatchScore = 1.5

	// Act
	err := jm.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "match score must be between 0 and 1")
}

func TestJobMatch_IsHighMatch_WhenScoreAbove70_ReturnsTrue(t *testing.T) {
	// Arrange
	profileID := uuid.New()
	jm := domain.NewJobMatch(profileID, "Looking for a Go developer...")
	jm.MatchScore = 0.85

	// Act
	result := jm.IsHighMatch()

	// Assert
	assert.True(t, result)
}

func TestJobMatch_IsMediumMatch_WhenScoreBetween40And70_ReturnsTrue(t *testing.T) {
	// Arrange
	profileID := uuid.New()
	jm := domain.NewJobMatch(profileID, "Looking for a Go developer...")
	jm.MatchScore = 0.55

	// Act
	result := jm.IsMediumMatch()

	// Assert
	assert.True(t, result)
}

func TestJobMatch_IsLowMatch_WhenScoreBelow40_ReturnsTrue(t *testing.T) {
	// Arrange
	profileID := uuid.New()
	jm := domain.NewJobMatch(profileID, "Looking for a Go developer...")
	jm.MatchScore = 0.2

	// Act
	result := jm.IsLowMatch()

	// Assert
	assert.True(t, result)
}

func TestJobMatch_SetMatchResult_UpdatesScoreAndSummary(t *testing.T) {
	// Arrange
	profileID := uuid.New()
	jm := domain.NewJobMatch(profileID, "Looking for a Go developer...")

	// Act
	jm.SetMatchResult(0.75, "Strong match for backend development")

	// Assert
	assert.Equal(t, 0.75, jm.MatchScore)
	assert.Equal(t, "Strong match for backend development", jm.TailoredSummary)
}
