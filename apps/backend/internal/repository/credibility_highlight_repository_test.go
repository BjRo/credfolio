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

// T048: Unit test for CredibilityHighlightRepository when saving highlight persists to database
func TestGormCredibilityHighlightRepository_Create_WhenSavingHighlight_PersistsToDatabase(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormCredibilityHighlightRepository(db)
	profileRepo := repository.NewGormProfileRepository(db)
	workExpRepo := repository.NewGormWorkExperienceRepository(db)
	refLetterRepo := repository.NewGormReferenceLetterRepository(db)
	ctx := context.Background()

	// Create prerequisites
	userID := uuid.New()
	profile := domain.NewProfile(userID)
	err := profileRepo.Create(ctx, profile)
	require.NoError(t, err)

	letter := domain.NewReferenceLetter(userID, "test.txt", "/uploads/test.txt")
	err = refLetterRepo.Create(ctx, letter)
	require.NoError(t, err)

	experience := &domain.WorkExperience{
		ID:          uuid.New(),
		ProfileID:   profile.ID,
		CompanyName: "Test Company",
		Role:        "Engineer",
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = workExpRepo.Create(ctx, experience)
	require.NoError(t, err)

	highlight := domain.NewCredibilityHighlight(
		"John was an exceptional team player",
		domain.SentimentPositive,
		letter.ID,
	)
	highlight.WorkExperienceID = experience.ID

	// Act
	err = repo.Create(ctx, highlight)

	// Assert
	require.NoError(t, err)

	// Verify persistence
	found, err := repo.FindByID(ctx, highlight.ID)
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "John was an exceptional team player", found.Quote)
	assert.Equal(t, domain.SentimentPositive, found.Sentiment)
}

func TestGormCredibilityHighlightRepository_FindByWorkExperienceID_WhenHighlightsExist_ReturnsHighlights(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormCredibilityHighlightRepository(db)
	profileRepo := repository.NewGormProfileRepository(db)
	workExpRepo := repository.NewGormWorkExperienceRepository(db)
	refLetterRepo := repository.NewGormReferenceLetterRepository(db)
	ctx := context.Background()

	// Create prerequisites
	userID := uuid.New()
	profile := domain.NewProfile(userID)
	err := profileRepo.Create(ctx, profile)
	require.NoError(t, err)

	letter := domain.NewReferenceLetter(userID, "test.txt", "/uploads/test.txt")
	err = refLetterRepo.Create(ctx, letter)
	require.NoError(t, err)

	experience := &domain.WorkExperience{
		ID:          uuid.New(),
		ProfileID:   profile.ID,
		CompanyName: "Test Company",
		Role:        "Engineer",
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = workExpRepo.Create(ctx, experience)
	require.NoError(t, err)

	highlight1 := domain.NewCredibilityHighlight("Quote 1", domain.SentimentPositive, letter.ID)
	highlight1.WorkExperienceID = experience.ID
	err = repo.Create(ctx, highlight1)
	require.NoError(t, err)

	highlight2 := domain.NewCredibilityHighlight("Quote 2", domain.SentimentNeutral, letter.ID)
	highlight2.WorkExperienceID = experience.ID
	err = repo.Create(ctx, highlight2)
	require.NoError(t, err)

	// Act
	highlights, err := repo.FindByWorkExperienceID(ctx, experience.ID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, highlights, 2)
}

func TestGormCredibilityHighlightRepository_FindBySourceLetterID_WhenHighlightsExist_ReturnsHighlights(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormCredibilityHighlightRepository(db)
	profileRepo := repository.NewGormProfileRepository(db)
	workExpRepo := repository.NewGormWorkExperienceRepository(db)
	refLetterRepo := repository.NewGormReferenceLetterRepository(db)
	ctx := context.Background()

	// Create prerequisites
	userID := uuid.New()
	profile := domain.NewProfile(userID)
	err := profileRepo.Create(ctx, profile)
	require.NoError(t, err)

	letter := domain.NewReferenceLetter(userID, "test.txt", "/uploads/test.txt")
	err = refLetterRepo.Create(ctx, letter)
	require.NoError(t, err)

	experience := &domain.WorkExperience{
		ID:          uuid.New(),
		ProfileID:   profile.ID,
		CompanyName: "Test Company",
		Role:        "Engineer",
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = workExpRepo.Create(ctx, experience)
	require.NoError(t, err)

	highlight := domain.NewCredibilityHighlight("Excellent work", domain.SentimentPositive, letter.ID)
	highlight.WorkExperienceID = experience.ID
	err = repo.Create(ctx, highlight)
	require.NoError(t, err)

	// Act
	highlights, err := repo.FindBySourceLetterID(ctx, letter.ID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, highlights, 1)
	assert.Equal(t, "Excellent work", highlights[0].Quote)
}

func TestGormCredibilityHighlightRepository_Delete_WhenHighlightExists_RemovesHighlight(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGormCredibilityHighlightRepository(db)
	profileRepo := repository.NewGormProfileRepository(db)
	workExpRepo := repository.NewGormWorkExperienceRepository(db)
	refLetterRepo := repository.NewGormReferenceLetterRepository(db)
	ctx := context.Background()

	// Create prerequisites
	userID := uuid.New()
	profile := domain.NewProfile(userID)
	err := profileRepo.Create(ctx, profile)
	require.NoError(t, err)

	letter := domain.NewReferenceLetter(userID, "test.txt", "/uploads/test.txt")
	err = refLetterRepo.Create(ctx, letter)
	require.NoError(t, err)

	experience := &domain.WorkExperience{
		ID:          uuid.New(),
		ProfileID:   profile.ID,
		CompanyName: "Test Company",
		Role:        "Engineer",
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = workExpRepo.Create(ctx, experience)
	require.NoError(t, err)

	highlight := domain.NewCredibilityHighlight("To delete", domain.SentimentPositive, letter.ID)
	highlight.WorkExperienceID = experience.ID
	err = repo.Create(ctx, highlight)
	require.NoError(t, err)

	// Act
	err = repo.Delete(ctx, highlight.ID)

	// Assert
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, highlight.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}
