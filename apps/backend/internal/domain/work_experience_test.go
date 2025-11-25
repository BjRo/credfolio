package domain_test

import (
	"testing"
	"time"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkExperience_Validate_WhenValidData_ReturnsNoError(t *testing.T) {
	// Arrange
	we := domain.NewWorkExperience("Test Company", "Software Engineer")

	// Act
	err := we.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestWorkExperience_Validate_WhenCompanyNameEmpty_ReturnsError(t *testing.T) {
	// Arrange
	we := domain.NewWorkExperience("", "Software Engineer")

	// Act
	err := we.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "company name is required")
}

func TestWorkExperience_Validate_WhenRoleEmpty_ReturnsError(t *testing.T) {
	// Arrange
	we := domain.NewWorkExperience("Test Company", "")

	// Act
	err := we.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "role is required")
}

func TestWorkExperience_Validate_WhenEndDateBeforeStartDate_ReturnsError(t *testing.T) {
	// Arrange
	we := domain.NewWorkExperience("Test Company", "Software Engineer")
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	we.SetDates(startDate, &endDate)

	// Act
	err := we.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "end date must be after start date")
}

func TestWorkExperience_IsCurrent_WhenNoEndDate_ReturnsTrue(t *testing.T) {
	// Arrange
	we := domain.NewWorkExperience("Test Company", "Software Engineer")
	we.EndDate = nil

	// Act
	result := we.IsCurrent()

	// Assert
	assert.True(t, result)
}

func TestWorkExperience_IsCurrent_WhenHasEndDate_ReturnsFalse(t *testing.T) {
	// Arrange
	we := domain.NewWorkExperience("Test Company", "Software Engineer")
	endDate := time.Now()
	we.EndDate = &endDate

	// Act
	result := we.IsCurrent()

	// Assert
	assert.False(t, result)
}

func TestNewWorkExperience_GeneratesUUID(t *testing.T) {
	// Act
	we := domain.NewWorkExperience("Test Company", "Software Engineer")

	// Assert
	assert.NotEqual(t, we.ID.String(), "00000000-0000-0000-0000-000000000000")
}
