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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupJobMatchTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Profile{},
		&domain.Skill{},
		&domain.WorkExperience{},
		&domain.CredibilityHighlight{},
		&domain.JobMatch{},
	)
	require.NoError(t, err)

	return db
}

// T103: Unit test for JobMatchRepository when saving job match persists to database
func TestJobMatchRepository_Create_WhenSavingJobMatch_PersistsToDatabase(t *testing.T) {
	// Arrange
	db := setupJobMatchTestDB(t)
	repo := repository.NewGormJobMatchRepository(db)
	ctx := context.Background()

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(user).Error
	require.NoError(t, err)

	profile := domain.NewProfile(user.ID)
	err = db.Create(profile).Error
	require.NoError(t, err)

	jobMatch := &domain.JobMatch{
		ID:              uuid.New(),
		BaseProfileID:   profile.ID,
		JobDescription:  "Software engineer with Go and TypeScript experience",
		MatchScore:      0.85,
		TailoredSummary: "Strong match based on technical skills",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Act
	err = repo.Create(ctx, jobMatch)

	// Assert
	require.NoError(t, err)

	var savedMatch domain.JobMatch
	err = db.First(&savedMatch, "id = ?", jobMatch.ID).Error
	require.NoError(t, err)
	assert.Equal(t, jobMatch.ID, savedMatch.ID)
	assert.Equal(t, float64(0.85), savedMatch.MatchScore)
}

// T104: Unit test for JobMatchRepository when finding by profile ID returns matches
func TestJobMatchRepository_FindByProfileID_WhenProfileHasMatches_ReturnsMatches(t *testing.T) {
	// Arrange
	db := setupJobMatchTestDB(t)
	repo := repository.NewGormJobMatchRepository(db)
	ctx := context.Background()

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(user).Error
	require.NoError(t, err)

	profile := domain.NewProfile(user.ID)
	err = db.Create(profile).Error
	require.NoError(t, err)

	jobMatches := []*domain.JobMatch{
		{
			ID:              uuid.New(),
			BaseProfileID:   profile.ID,
			JobDescription:  "Frontend engineer",
			MatchScore:      0.75,
			TailoredSummary: "Good match",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New(),
			BaseProfileID:   profile.ID,
			JobDescription:  "Backend engineer",
			MatchScore:      0.90,
			TailoredSummary: "Excellent match",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	for _, jm := range jobMatches {
		err = db.Create(jm).Error
		require.NoError(t, err)
	}

	// Act
	results, err := repo.FindByProfileID(ctx, profile.ID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestJobMatchRepository_FindByID_WhenJobMatchExists_ReturnsMatch(t *testing.T) {
	// Arrange
	db := setupJobMatchTestDB(t)
	repo := repository.NewGormJobMatchRepository(db)
	ctx := context.Background()

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(user).Error
	require.NoError(t, err)

	profile := domain.NewProfile(user.ID)
	err = db.Create(profile).Error
	require.NoError(t, err)

	jobMatch := &domain.JobMatch{
		ID:              uuid.New(),
		BaseProfileID:   profile.ID,
		JobDescription:  "Software engineer",
		MatchScore:      0.80,
		TailoredSummary: "Good match",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = db.Create(jobMatch).Error
	require.NoError(t, err)

	// Act
	result, err := repo.FindByID(ctx, jobMatch.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, jobMatch.ID, result.ID)
	assert.Equal(t, float64(0.80), result.MatchScore)
}

func TestJobMatchRepository_FindByID_WhenJobMatchNotExists_ReturnsNil(t *testing.T) {
	// Arrange
	db := setupJobMatchTestDB(t)
	repo := repository.NewGormJobMatchRepository(db)
	ctx := context.Background()

	// Act
	result, err := repo.FindByID(ctx, uuid.New())

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result)
}
