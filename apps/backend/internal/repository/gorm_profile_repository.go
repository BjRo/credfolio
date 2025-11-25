package repository

import (
	"context"
	"errors"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormProfileRepository implements ProfileRepository using GORM
type GormProfileRepository struct {
	db *gorm.DB
}

// NewGormProfileRepository creates a new GORM profile repository
func NewGormProfileRepository(db *gorm.DB) *GormProfileRepository {
	return &GormProfileRepository{db: db}
}

// Create creates a new profile
func (r *GormProfileRepository) Create(ctx context.Context, profile *domain.Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// Update updates an existing profile
func (r *GormProfileRepository) Update(ctx context.Context, profile *domain.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// FindByID finds a profile by its ID
func (r *GormProfileRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.db.WithContext(ctx).First(&profile, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// FindByUserID finds a profile by user ID
func (r *GormProfileRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.db.WithContext(ctx).First(&profile, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// Delete deletes a profile
func (r *GormProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Profile{}, "id = ?", id).Error
}

// FindByUserIDWithRelations finds a profile with all relations loaded
func (r *GormProfileRepository) FindByUserIDWithRelations(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.db.WithContext(ctx).
		Preload("WorkExperiences").
		Preload("WorkExperiences.CredibilityHighlights").
		Preload("Skills").
		First(&profile, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// Ensure GormProfileRepository implements ProfileRepository
var _ ProfileRepository = (*GormProfileRepository)(nil)
