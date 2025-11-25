package repository

import (
	"github.com/credfolio/apps/backend/internal/domain"
)

// RunMigrations runs all database migrations using GORM AutoMigrate
func RunMigrations(db *Database) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Profile{},
		&domain.WorkExperience{},
		&domain.Skill{},
		&domain.ProfileSkill{},
		&domain.ReferenceLetter{},
		&domain.CredibilityHighlight{},
		&domain.JobMatch{},
	)
}
