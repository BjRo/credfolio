package domain_test

import (
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfile_Validate_WhenValidData_ReturnsNoError(t *testing.T) {
	// Arrange
	userID := uuid.New()
	profile := domain.NewProfile(userID)

	// Act
	err := profile.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestProfile_Validate_WhenUserIDEmpty_ReturnsError(t *testing.T) {
	// Arrange
	profile := domain.NewProfile(uuid.Nil)

	// Act
	err := profile.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user ID is required")
}

func TestNewProfile_GeneratesUUID(t *testing.T) {
	// Arrange
	userID := uuid.New()

	// Act
	profile := domain.NewProfile(userID)

	// Assert
	assert.NotEqual(t, profile.ID.String(), "00000000-0000-0000-0000-000000000000")
	assert.Equal(t, userID, profile.UserID)
}

func TestProfile_AddWorkExperience_SetsProfileID(t *testing.T) {
	// Arrange
	userID := uuid.New()
	profile := domain.NewProfile(userID)
	we := domain.NewWorkExperience("Test Company", "Engineer")

	// Act
	profile.AddWorkExperience(we)

	// Assert
	assert.Len(t, profile.WorkExperiences, 1)
	assert.Equal(t, profile.ID, profile.WorkExperiences[0].ProfileID)
}

func TestProfile_AddSkill_AddsToSkillsList(t *testing.T) {
	// Arrange
	userID := uuid.New()
	profile := domain.NewProfile(userID)
	skill := domain.NewSkill("Go")

	// Act
	profile.AddSkill(skill)

	// Assert
	assert.Len(t, profile.Skills, 1)
	assert.Equal(t, "Go", profile.Skills[0].Name)
}
