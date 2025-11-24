package repository

import (
	"context"
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGormReferenceLetterRepository_Create_WhenValidLetter_ThenSucceeds(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewGormReferenceLetterRepository(db)
	ctx := context.Background()
	userID := uuid.New()
	letter := &domain.ReferenceLetter{
		UserID:      userID,
		FileName:    "test.pdf",
		StoragePath: "/storage/test.pdf",
		Status:      domain.ReferenceLetterStatusPending,
	}

	// Act
	err := repo.Create(ctx, letter)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, letter.ID)
}

func TestGormReferenceLetterRepository_GetByUserID_WhenLettersExist_ThenReturnsAllLetters(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewGormReferenceLetterRepository(db)
	ctx := context.Background()
	userID := uuid.New()

	letter1 := &domain.ReferenceLetter{
		UserID:      userID,
		FileName:    "test1.pdf",
		StoragePath: "/storage/test1.pdf",
		Status:      domain.ReferenceLetterStatusPending,
	}
	letter2 := &domain.ReferenceLetter{
		UserID:      userID,
		FileName:    "test2.pdf",
		StoragePath: "/storage/test2.pdf",
		Status:      domain.ReferenceLetterStatusPending,
	}
	err := repo.Create(ctx, letter1)
	require.NoError(t, err)
	err = repo.Create(ctx, letter2)
	require.NoError(t, err)

	// Act
	letters, err := repo.GetByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, letters, 2)
}

func TestGormReferenceLetterRepository_Update_WhenValidLetter_ThenSucceeds(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewGormReferenceLetterRepository(db)
	ctx := context.Background()
	userID := uuid.New()
	letter := &domain.ReferenceLetter{
		UserID:      userID,
		FileName:    "test.pdf",
		StoragePath: "/storage/test.pdf",
		Status:      domain.ReferenceLetterStatusPending,
	}
	err := repo.Create(ctx, letter)
	require.NoError(t, err)

	letter.Status = domain.ReferenceLetterStatusProcessed
	letter.ExtractedText = "Extracted text content"

	// Act
	err = repo.Update(ctx, letter)

	// Assert
	assert.NoError(t, err)
	updated, err := repo.GetByID(ctx, letter.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.ReferenceLetterStatusProcessed, updated.Status)
	assert.Equal(t, "Extracted text content", updated.ExtractedText)
}

func TestGormReferenceLetterRepository_Delete_WhenLetterExists_ThenSucceeds(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewGormReferenceLetterRepository(db)
	ctx := context.Background()
	userID := uuid.New()
	letter := &domain.ReferenceLetter{
		UserID:      userID,
		FileName:    "test.pdf",
		StoragePath: "/storage/test.pdf",
		Status:      domain.ReferenceLetterStatusPending,
	}
	err := repo.Create(ctx, letter)
	require.NoError(t, err)

	// Act
	err = repo.Delete(ctx, letter.ID)

	// Assert
	assert.NoError(t, err)
	_, err = repo.GetByID(ctx, letter.ID)
	assert.Error(t, err)
}
