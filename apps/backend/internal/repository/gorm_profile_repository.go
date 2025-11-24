package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormProfileRepository implements ProfileRepository using GORM
type GormProfileRepository struct {
	db *gorm.DB
}

// NewGormProfileRepository creates a new GORM-based profile repository
func NewGormProfileRepository(db *gorm.DB) *GormProfileRepository {
	return &GormProfileRepository{db: db}
}

// Create creates a new profile
func (r *GormProfileRepository) Create(ctx context.Context, profile *domain.Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// GetByID retrieves a profile by ID
func (r *GormProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.db.WithContext(ctx).
		Preload("WorkExperiences").
		Preload("WorkExperiences.CredibilityHighlights").
		Preload("Skills").
		Preload("JobMatches").
		First(&profile, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetByUserID retrieves a profile by user ID
func (r *GormProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.db.WithContext(ctx).
		Preload("WorkExperiences").
		Preload("WorkExperiences.CredibilityHighlights").
		Preload("Skills").
		Preload("JobMatches").
		First(&profile, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Update updates an existing profile
func (r *GormProfileRepository) Update(ctx context.Context, profile *domain.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// Delete deletes a profile
func (r *GormProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Profile{}, "id = ?", id).Error
}
