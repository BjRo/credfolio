package repository

import (
	"context"
	"errors"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormWorkExperienceRepository implements WorkExperienceRepository using GORM
type GormWorkExperienceRepository struct {
	db *gorm.DB
}

// NewGormWorkExperienceRepository creates a new GORM work experience repository
func NewGormWorkExperienceRepository(db *gorm.DB) *GormWorkExperienceRepository {
	return &GormWorkExperienceRepository{db: db}
}

// Create creates a new work experience
func (r *GormWorkExperienceRepository) Create(ctx context.Context, experience *domain.WorkExperience) error {
	return r.db.WithContext(ctx).Create(experience).Error
}

// Update updates an existing work experience
func (r *GormWorkExperienceRepository) Update(ctx context.Context, experience *domain.WorkExperience) error {
	return r.db.WithContext(ctx).Save(experience).Error
}

// FindByID finds a work experience by its ID
func (r *GormWorkExperienceRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.WorkExperience, error) {
	var experience domain.WorkExperience
	err := r.db.WithContext(ctx).First(&experience, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &experience, nil
}

// FindByProfileID finds all work experiences for a profile
func (r *GormWorkExperienceRepository) FindByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.WorkExperience, error) {
	var experiences []*domain.WorkExperience
	err := r.db.WithContext(ctx).Find(&experiences, "profile_id = ?", profileID).Error
	if err != nil {
		return nil, err
	}
	return experiences, nil
}

// FindByProfileIDWithHighlights finds work experiences with credibility highlights
func (r *GormWorkExperienceRepository) FindByProfileIDWithHighlights(ctx context.Context, profileID uuid.UUID) ([]*domain.WorkExperience, error) {
	var experiences []*domain.WorkExperience
	err := r.db.WithContext(ctx).
		Preload("CredibilityHighlights").
		Find(&experiences, "profile_id = ?", profileID).Error
	if err != nil {
		return nil, err
	}
	return experiences, nil
}

// Delete deletes a work experience
func (r *GormWorkExperienceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.WorkExperience{}, "id = ?", id).Error
}

// Ensure GormWorkExperienceRepository implements WorkExperienceRepository
var _ WorkExperienceRepository = (*GormWorkExperienceRepository)(nil)
