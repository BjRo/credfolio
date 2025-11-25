package repository

import (
	"context"
	"errors"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormCredibilityHighlightRepository implements CredibilityHighlightRepository using GORM
type GormCredibilityHighlightRepository struct {
	db *gorm.DB
}

// NewGormCredibilityHighlightRepository creates a new GORM credibility highlight repository
func NewGormCredibilityHighlightRepository(db *gorm.DB) *GormCredibilityHighlightRepository {
	return &GormCredibilityHighlightRepository{db: db}
}

// Create creates a new credibility highlight
func (r *GormCredibilityHighlightRepository) Create(ctx context.Context, highlight *domain.CredibilityHighlight) error {
	return r.db.WithContext(ctx).Create(highlight).Error
}

// Update updates an existing credibility highlight
func (r *GormCredibilityHighlightRepository) Update(ctx context.Context, highlight *domain.CredibilityHighlight) error {
	return r.db.WithContext(ctx).Save(highlight).Error
}

// FindByID finds a credibility highlight by its ID
func (r *GormCredibilityHighlightRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.CredibilityHighlight, error) {
	var highlight domain.CredibilityHighlight
	err := r.db.WithContext(ctx).First(&highlight, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &highlight, nil
}

// FindByWorkExperienceID finds all highlights for a work experience
func (r *GormCredibilityHighlightRepository) FindByWorkExperienceID(ctx context.Context, workExperienceID uuid.UUID) ([]*domain.CredibilityHighlight, error) {
	var highlights []*domain.CredibilityHighlight
	err := r.db.WithContext(ctx).Find(&highlights, "work_experience_id = ?", workExperienceID).Error
	if err != nil {
		return nil, err
	}
	return highlights, nil
}

// FindBySourceLetterID finds all highlights from a reference letter
func (r *GormCredibilityHighlightRepository) FindBySourceLetterID(ctx context.Context, sourceLetterID uuid.UUID) ([]*domain.CredibilityHighlight, error) {
	var highlights []*domain.CredibilityHighlight
	err := r.db.WithContext(ctx).Find(&highlights, "source_letter_id = ?", sourceLetterID).Error
	if err != nil {
		return nil, err
	}
	return highlights, nil
}

// Delete deletes a credibility highlight
func (r *GormCredibilityHighlightRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.CredibilityHighlight{}, "id = ?", id).Error
}

// Ensure GormCredibilityHighlightRepository implements CredibilityHighlightRepository
var _ CredibilityHighlightRepository = (*GormCredibilityHighlightRepository)(nil)
