package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormCredibilityHighlightRepository implements CredibilityHighlightRepository using GORM
type GormCredibilityHighlightRepository struct {
	db *gorm.DB
}

// NewGormCredibilityHighlightRepository creates a new GORM-based credibility highlight repository
func NewGormCredibilityHighlightRepository(db *gorm.DB) *GormCredibilityHighlightRepository {
	return &GormCredibilityHighlightRepository{db: db}
}

// Create creates a new credibility highlight
func (r *GormCredibilityHighlightRepository) Create(ctx context.Context, highlight *domain.CredibilityHighlight) error {
	return r.db.WithContext(ctx).Create(highlight).Error
}

// GetByID retrieves a credibility highlight by ID
func (r *GormCredibilityHighlightRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CredibilityHighlight, error) {
	var highlight domain.CredibilityHighlight
	err := r.db.WithContext(ctx).First(&highlight, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &highlight, nil
}

// GetByWorkExperienceID retrieves all credibility highlights for a work experience
func (r *GormCredibilityHighlightRepository) GetByWorkExperienceID(ctx context.Context, workExpID uuid.UUID) ([]*domain.CredibilityHighlight, error) {
	var highlights []*domain.CredibilityHighlight
	err := r.db.WithContext(ctx).Where("work_experience_id = ?", workExpID).Find(&highlights).Error
	return highlights, err
}

// Update updates an existing credibility highlight
func (r *GormCredibilityHighlightRepository) Update(ctx context.Context, highlight *domain.CredibilityHighlight) error {
	return r.db.WithContext(ctx).Save(highlight).Error
}

// Delete deletes a credibility highlight
func (r *GormCredibilityHighlightRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.CredibilityHighlight{}, "id = ?", id).Error
}
