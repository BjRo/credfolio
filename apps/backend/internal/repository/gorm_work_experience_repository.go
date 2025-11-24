package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormWorkExperienceRepository implements WorkExperienceRepository using GORM
type GormWorkExperienceRepository struct {
	db *gorm.DB
}

// NewGormWorkExperienceRepository creates a new GORM-based work experience repository
func NewGormWorkExperienceRepository(db *gorm.DB) *GormWorkExperienceRepository {
	return &GormWorkExperienceRepository{db: db}
}

// Create creates a new work experience
func (r *GormWorkExperienceRepository) Create(ctx context.Context, exp *domain.WorkExperience) error {
	return r.db.WithContext(ctx).Create(exp).Error
}

// GetByID retrieves a work experience by ID
func (r *GormWorkExperienceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WorkExperience, error) {
	var exp domain.WorkExperience
	err := r.db.WithContext(ctx).
		Preload("CredibilityHighlights").
		First(&exp, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

// GetByProfileID retrieves all work experiences for a profile
func (r *GormWorkExperienceRepository) GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.WorkExperience, error) {
	var experiences []*domain.WorkExperience
	err := r.db.WithContext(ctx).
		Preload("CredibilityHighlights").
		Where("profile_id = ?", profileID).
		Find(&experiences).Error
	return experiences, err
}

// Update updates an existing work experience
func (r *GormWorkExperienceRepository) Update(ctx context.Context, exp *domain.WorkExperience) error {
	return r.db.WithContext(ctx).Save(exp).Error
}

// Delete deletes a work experience
func (r *GormWorkExperienceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.WorkExperience{}, "id = ?", id).Error
}
