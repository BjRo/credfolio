package repository

import (
	"fmt"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// RunMigrations runs GORM AutoMigrate for all domain models
func RunMigrations() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	err := DB.AutoMigrate(
		&domain.User{},
		&domain.Profile{},
		&domain.WorkExperience{},
		&domain.Skill{},
		&domain.ProfileSkill{},
		&domain.ReferenceLetter{},
		&domain.CredibilityHighlight{},
		&domain.JobMatch{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Ensure the mock user exists for development
	if err := EnsureMockUser(); err != nil {
		return fmt.Errorf("failed to ensure mock user: %w", err)
	}

	return nil
}

// EnsureMockUser ensures the mock user exists in the database
// This is used for development when there's no real authentication
// The mock user ID matches the one in middleware/auth.go
func EnsureMockUser() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Mock user ID matches the one in middleware/auth.go
	mockUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	mockUser := &domain.User{
		ID:    mockUserID,
		Email: "dev@credfolio.local",
		Name:  "Development User",
	}

	// FirstOrCreate will create the user if it doesn't exist, or do nothing if it does
	result := DB.Where("id = ?", mockUserID).FirstOrCreate(mockUser)
	if result.Error != nil {
		return fmt.Errorf("failed to ensure mock user: %w", result.Error)
	}

	return nil
}
