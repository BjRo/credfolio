package repository_test

import (
	"context"
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T045: Unit test for ReferenceLetterRepository when saving reference letter persists to database
func TestGormReferenceLetterRepository_Create_WhenSavingReferenceLetter_PersistsToDatabase(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormReferenceLetterRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	letter := domain.NewReferenceLetter(userID, "test.txt", "/uploads/test.txt")

	// Act
	err := repo.Create(ctx, letter)

	// Assert
	require.NoError(t, err)

	// Verify persistence
	found, err := repo.FindByID(ctx, letter.ID)
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, letter.ID, found.ID)
	assert.Equal(t, userID, found.UserID)
	assert.Equal(t, "test.txt", found.FileName)
}

// T046: Unit test for ReferenceLetterRepository when finding by user ID returns letters
func TestGormReferenceLetterRepository_FindByUserID_WhenLettersExist_ReturnsLetters(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormReferenceLetterRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	letter1 := domain.NewReferenceLetter(userID, "letter1.txt", "/uploads/letter1.txt")
	letter2 := domain.NewReferenceLetter(userID, "letter2.txt", "/uploads/letter2.txt")

	err := repo.Create(ctx, letter1)
	require.NoError(t, err)
	err = repo.Create(ctx, letter2)
	require.NoError(t, err)

	// Act
	letters, err := repo.FindByUserID(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, letters, 2)

	fileNames := []string{letters[0].FileName, letters[1].FileName}
	assert.Contains(t, fileNames, "letter1.txt")
	assert.Contains(t, fileNames, "letter2.txt")
}

func TestGormReferenceLetterRepository_FindByUserID_WhenNoLettersExist_ReturnsEmptySlice(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormReferenceLetterRepository(db)
	ctx := context.Background()

	nonExistentUserID := uuid.New()

	// Act
	letters, err := repo.FindByUserID(ctx, nonExistentUserID)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, letters)
}

func TestGormReferenceLetterRepository_FindPendingByUserID_WhenPendingLettersExist_ReturnsPendingOnly(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormReferenceLetterRepository(db)
	ctx := context.Background()

	userID := uuid.New()

	pendingLetter := domain.NewReferenceLetter(userID, "pending.txt", "/uploads/pending.txt")
	pendingLetter.ExtractedText = "Some extracted text"

	processedLetter := domain.NewReferenceLetter(userID, "processed.txt", "/uploads/processed.txt")
	processedLetter.MarkProcessed("Processed text")

	err := repo.Create(ctx, pendingLetter)
	require.NoError(t, err)
	err = repo.Create(ctx, processedLetter)
	require.NoError(t, err)

	// Act
	pendingLetters, err := repo.FindPendingByUserID(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, pendingLetters, 1)
	assert.Equal(t, "pending.txt", pendingLetters[0].FileName)
	assert.Equal(t, domain.StatusPending, pendingLetters[0].Status)
}

func TestGormReferenceLetterRepository_Update_WhenLetterExists_UpdatesLetter(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormReferenceLetterRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	letter := domain.NewReferenceLetter(userID, "test.txt", "/uploads/test.txt")

	err := repo.Create(ctx, letter)
	require.NoError(t, err)

	// Act
	letter.MarkProcessed("Extracted content")
	err = repo.Update(ctx, letter)

	// Assert
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, letter.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.StatusProcessed, found.Status)
	assert.Equal(t, "Extracted content", found.ExtractedText)
}

func TestGormReferenceLetterRepository_Delete_WhenLetterExists_RemovesLetter(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormReferenceLetterRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	letter := domain.NewReferenceLetter(userID, "test.txt", "/uploads/test.txt")

	err := repo.Create(ctx, letter)
	require.NoError(t, err)

	// Act
	err = repo.Delete(ctx, letter.ID)

	// Assert
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, letter.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}
