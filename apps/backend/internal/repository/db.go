package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/credfolio/apps/backend/pkg/config"
)

// DB is the global database connection
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(cfg *config.Config) error {
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}
