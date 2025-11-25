package domain_test

import (
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkill_Validate_WhenValidData_ReturnsNoError(t *testing.T) {
	// Arrange
	skill := domain.NewSkill("Go")

	// Act
	err := skill.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestSkill_Validate_WhenNameEmpty_ReturnsError(t *testing.T) {
	// Arrange
	skill := domain.NewSkill("")

	// Act
	err := skill.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "skill name is required")
}

func TestNewSkill_GeneratesUUID(t *testing.T) {
	// Act
	skill := domain.NewSkill("Go")

	// Assert
	assert.NotEqual(t, skill.ID.String(), "00000000-0000-0000-0000-000000000000")
	assert.Equal(t, "Go", skill.Name)
}
