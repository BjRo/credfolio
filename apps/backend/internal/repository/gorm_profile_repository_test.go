package repository

import (
	"context"
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Profile{},
		&domain.WorkExperience{},
		&domain.Skill{},
		&domain.ProfileSkill{},
		&domain.ReferenceLetter{},
		&domain.CredibilityHighlight{},
		&domain.JobMatch{},
	)
	require.NoError(t, err)

	return db
}

func TestGormProfileRepository_Create_WhenValidProfile_ThenSucceeds(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewGormProfileRepository(db)
	ctx := context.Background()
	userID := uuid.New()
	profile := &domain.Profile{
		UserID:  userID,
		Summary: "Test summary",
	}

	// Act
	err := repo.Create(ctx, profile)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, profile.ID)
}

func TestGormProfileRepository_GetByID_WhenProfileExists_ThenReturnsProfile(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewGormProfileRepository(db)
	ctx := context.Background()
	userID := uuid.New()
	profile := &domain.Profile{
		UserID:  userID,
		Summary: "Test summary",
	}
	err := repo.Create(ctx, profile)
	require.NoError(t, err)

	// Act
	retrieved, err := repo.GetByID(ctx, profile.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, profile.ID, retrieved.ID)
	assert.Equal(t, userID, retrieved.UserID)
	assert.Equal(t, "Test summary", retrieved.Summary)
}

func TestGormProfileRepository_GetByID_WhenProfileNotFound_ThenReturnsError(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewGormProfileRepository(db)
	ctx := context.Background()
	nonExistentID := uuid.New()

	// Act
	retrieved, err := repo.GetByID(ctx, nonExistentID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrieved)
}

func TestGormProfileRepository_GetByUserID_WhenProfileExists_ThenReturnsProfile(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewGormProfileRepository(db)
	ctx := context.Background()
	userID := uuid.New()
	profile := &domain.Profile{
		UserID:  userID,
		Summary: "Test summary",
	}
	err := repo.Create(ctx, profile)
	require.NoError(t, err)

	// Act
	retrieved, err := repo.GetByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, profile.ID, retrieved.ID)
	assert.Equal(t, userID, retrieved.UserID)
}

func TestGormProfileRepository_Update_WhenValidProfile_ThenSucceeds(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewGormProfileRepository(db)
	ctx := context.Background()
	userID := uuid.New()
	profile := &domain.Profile{
		UserID:  userID,
		Summary: "Original summary",
	}
	err := repo.Create(ctx, profile)
	require.NoError(t, err)

	profile.Summary = "Updated summary"

	// Act
	err = repo.Update(ctx, profile)

	// Assert
	assert.NoError(t, err)
	updated, err := repo.GetByID(ctx, profile.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated summary", updated.Summary)
}

func TestGormProfileRepository_Delete_WhenProfileExists_ThenSucceeds(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewGormProfileRepository(db)
	ctx := context.Background()
	userID := uuid.New()
	profile := &domain.Profile{
		UserID:  userID,
		Summary: "Test summary",
	}
	err := repo.Create(ctx, profile)
	require.NoError(t, err)

	// Act
	err = repo.Delete(ctx, profile.ID)

	// Assert
	assert.NoError(t, err)
	_, err = repo.GetByID(ctx, profile.ID)
	assert.Error(t, err)
}
