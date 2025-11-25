package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/credfolio/apps/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T047: Unit test for WorkExperienceRepository when saving work experience persists to database
func TestGormWorkExperienceRepository_Create_WhenSavingWorkExperience_PersistsToDatabase(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormWorkExperienceRepository(db)
	profileRepo := repository.NewGormProfileRepository(db)
	ctx := context.Background()

	// Create profile first
	userID := uuid.New()
	profile := domain.NewProfile(userID)
	err := profileRepo.Create(ctx, profile)
	require.NoError(t, err)

	experience := &domain.WorkExperience{
		ID:          uuid.New(),
		ProfileID:   profile.ID,
		CompanyName: "Acme Corp",
		Role:        "Software Engineer",
		Description: "Built amazing software",
		StartDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Act
	err = repo.Create(ctx, experience)

	// Assert
	require.NoError(t, err)

	// Verify persistence
	found, err := repo.FindByID(ctx, experience.ID)
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, experience.ID, found.ID)
	assert.Equal(t, "Acme Corp", found.CompanyName)
	assert.Equal(t, "Software Engineer", found.Role)
}

func TestGormWorkExperienceRepository_FindByProfileID_WhenExperiencesExist_ReturnsExperiences(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormWorkExperienceRepository(db)
	profileRepo := repository.NewGormProfileRepository(db)
	ctx := context.Background()

	// Create profile first
	userID := uuid.New()
	profile := domain.NewProfile(userID)
	err := profileRepo.Create(ctx, profile)
	require.NoError(t, err)

	exp1 := &domain.WorkExperience{
		ID:          uuid.New(),
		ProfileID:   profile.ID,
		CompanyName: "Company A",
		Role:        "Engineer",
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	exp2 := &domain.WorkExperience{
		ID:          uuid.New(),
		ProfileID:   profile.ID,
		CompanyName: "Company B",
		Role:        "Senior Engineer",
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = repo.Create(ctx, exp1)
	require.NoError(t, err)
	err = repo.Create(ctx, exp2)
	require.NoError(t, err)

	// Act
	experiences, err := repo.FindByProfileID(ctx, profile.ID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, experiences, 2)
}

func TestGormWorkExperienceRepository_Update_WhenExperienceExists_UpdatesExperience(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormWorkExperienceRepository(db)
	profileRepo := repository.NewGormProfileRepository(db)
	ctx := context.Background()

	// Create profile first
	userID := uuid.New()
	profile := domain.NewProfile(userID)
	err := profileRepo.Create(ctx, profile)
	require.NoError(t, err)

	experience := &domain.WorkExperience{
		ID:          uuid.New(),
		ProfileID:   profile.ID,
		CompanyName: "Original Company",
		Role:        "Engineer",
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = repo.Create(ctx, experience)
	require.NoError(t, err)

	// Act
	experience.CompanyName = "Updated Company"
	err = repo.Update(ctx, experience)

	// Assert
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, experience.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Company", found.CompanyName)
}

func TestGormWorkExperienceRepository_Delete_WhenExperienceExists_RemovesExperience(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormWorkExperienceRepository(db)
	profileRepo := repository.NewGormProfileRepository(db)
	ctx := context.Background()

	// Create profile first
	userID := uuid.New()
	profile := domain.NewProfile(userID)
	err := profileRepo.Create(ctx, profile)
	require.NoError(t, err)

	experience := &domain.WorkExperience{
		ID:          uuid.New(),
		ProfileID:   profile.ID,
		CompanyName: "Company",
		Role:        "Engineer",
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = repo.Create(ctx, experience)
	require.NoError(t, err)

	// Act
	err = repo.Delete(ctx, experience.ID)

	// Assert
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, experience.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}
