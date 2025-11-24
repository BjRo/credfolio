package repository

import (
	"fmt"

	"github.com/credfolio/apps/backend/internal/domain"
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

	return nil
}
