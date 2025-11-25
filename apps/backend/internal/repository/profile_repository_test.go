package repository_test

import (
	"context"
	"testing"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Migrate schema
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Profile{},
		&domain.WorkExperience{},
		&domain.Skill{},
		&domain.ReferenceLetter{},
		&domain.CredibilityHighlight{},
	)
	require.NoError(t, err)

	return db
}

// T043: Unit test for ProfileRepository when saving profile persists to database
func TestGormProfileRepository_Create_WhenSavingProfile_PersistsToDatabase(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormProfileRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	profile := domain.NewProfile(userID)

	// Act
	err := repo.Create(ctx, profile)

	// Assert
	require.NoError(t, err)

	// Verify persistence
	found, err := repo.FindByID(ctx, profile.ID)
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, profile.ID, found.ID)
	assert.Equal(t, userID, found.UserID)
}

// T044: Unit test for ProfileRepository when finding by user ID returns profile
func TestGormProfileRepository_FindByUserID_WhenProfileExists_ReturnsProfile(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormProfileRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	profile := domain.NewProfile(userID)
	profile.Summary = "Test summary"

	err := repo.Create(ctx, profile)
	require.NoError(t, err)

	// Act
	found, err := repo.FindByUserID(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, userID, found.UserID)
	assert.Equal(t, "Test summary", found.Summary)
}

func TestGormProfileRepository_FindByUserID_WhenProfileNotExists_ReturnsNil(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormProfileRepository(db)
	ctx := context.Background()

	nonExistentUserID := uuid.New()

	// Act
	found, err := repo.FindByUserID(ctx, nonExistentUserID)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, found)
}

func TestGormProfileRepository_Update_WhenProfileExists_UpdatesProfile(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormProfileRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	profile := domain.NewProfile(userID)
	profile.Summary = "Original summary"
	err := repo.Create(ctx, profile)
	require.NoError(t, err)

	// Act
	profile.Summary = "Updated summary"
	err = repo.Update(ctx, profile)

	// Assert
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, profile.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated summary", found.Summary)
}

func TestGormProfileRepository_Delete_WhenProfileExists_RemovesProfile(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormProfileRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	profile := domain.NewProfile(userID)
	err := repo.Create(ctx, profile)
	require.NoError(t, err)

	// Act
	err = repo.Delete(ctx, profile.ID)

	// Assert
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, profile.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}
