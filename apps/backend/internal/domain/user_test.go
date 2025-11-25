package domain_test

import (
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Validate_WhenValidData_ReturnsNoError(t *testing.T) {
	// Arrange
	user := domain.NewUser("test@example.com", "John Doe")

	// Act
	err := user.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestUser_Validate_WhenEmailEmpty_ReturnsError(t *testing.T) {
	// Arrange
	user := domain.NewUser("", "John Doe")

	// Act
	err := user.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email is required")
}

func TestUser_Validate_WhenNameEmpty_ReturnsError(t *testing.T) {
	// Arrange
	user := domain.NewUser("test@example.com", "")

	// Act
	err := user.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestNewUser_GeneratesUUID(t *testing.T) {
	// Act
	user := domain.NewUser("test@example.com", "John Doe")

	// Assert
	assert.NotEqual(t, user.ID.String(), "00000000-0000-0000-0000-000000000000")
}
