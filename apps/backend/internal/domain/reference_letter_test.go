package domain_test

import (
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReferenceLetter_Validate_WhenValidData_ReturnsNoError(t *testing.T) {
	// Arrange
	userID := uuid.New()
	rl := domain.NewReferenceLetter(userID, "letter.pdf", "/uploads/letter.pdf")

	// Act
	err := rl.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestReferenceLetter_Validate_WhenUserIDEmpty_ReturnsError(t *testing.T) {
	// Arrange
	rl := domain.NewReferenceLetter(uuid.Nil, "letter.pdf", "/uploads/letter.pdf")

	// Act
	err := rl.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user ID is required")
}

func TestReferenceLetter_Validate_WhenFileNameEmpty_ReturnsError(t *testing.T) {
	// Arrange
	userID := uuid.New()
	rl := domain.NewReferenceLetter(userID, "", "/uploads/letter.pdf")

	// Act
	err := rl.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file name is required")
}

func TestReferenceLetter_Validate_WhenStoragePathEmpty_ReturnsError(t *testing.T) {
	// Arrange
	userID := uuid.New()
	rl := domain.NewReferenceLetter(userID, "letter.pdf", "")

	// Act
	err := rl.Validate()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "storage path is required")
}

func TestNewReferenceLetter_SetsStatusToPending(t *testing.T) {
	// Arrange
	userID := uuid.New()

	// Act
	rl := domain.NewReferenceLetter(userID, "letter.pdf", "/uploads/letter.pdf")

	// Assert
	assert.Equal(t, domain.StatusPending, rl.Status)
	assert.True(t, rl.IsPending())
}

func TestReferenceLetter_MarkProcessed_UpdatesStatusAndText(t *testing.T) {
	// Arrange
	userID := uuid.New()
	rl := domain.NewReferenceLetter(userID, "letter.pdf", "/uploads/letter.pdf")
	extractedText := "This is the extracted text from the letter."

	// Act
	rl.MarkProcessed(extractedText)

	// Assert
	assert.Equal(t, domain.StatusProcessed, rl.Status)
	assert.Equal(t, extractedText, rl.ExtractedText)
	assert.True(t, rl.IsProcessed())
}

func TestReferenceLetter_MarkFailed_UpdatesStatus(t *testing.T) {
	// Arrange
	userID := uuid.New()
	rl := domain.NewReferenceLetter(userID, "letter.pdf", "/uploads/letter.pdf")

	// Act
	rl.MarkFailed()

	// Assert
	assert.Equal(t, domain.StatusFailed, rl.Status)
	assert.False(t, rl.IsProcessed())
	assert.False(t, rl.IsPending())
}
