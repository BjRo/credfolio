package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormReferenceLetterRepository implements ReferenceLetterRepository using GORM
type GormReferenceLetterRepository struct {
	db *gorm.DB
}

// NewGormReferenceLetterRepository creates a new GORM-based reference letter repository
func NewGormReferenceLetterRepository(db *gorm.DB) *GormReferenceLetterRepository {
	return &GormReferenceLetterRepository{db: db}
}

// Create creates a new reference letter
func (r *GormReferenceLetterRepository) Create(ctx context.Context, letter *domain.ReferenceLetter) error {
	return r.db.WithContext(ctx).Create(letter).Error
}

// GetByID retrieves a reference letter by ID
func (r *GormReferenceLetterRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ReferenceLetter, error) {
	var letter domain.ReferenceLetter
	err := r.db.WithContext(ctx).First(&letter, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &letter, nil
}

// GetByUserID retrieves all reference letters for a user
func (r *GormReferenceLetterRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	var letters []*domain.ReferenceLetter
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&letters).Error
	return letters, err
}

// GetByContentSHA retrieves a reference letter by content SHA for a specific user
func (r *GormReferenceLetterRepository) GetByContentSHA(ctx context.Context, userID uuid.UUID, contentSHA string) (*domain.ReferenceLetter, error) {
	var letter domain.ReferenceLetter
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND content_sha = ?", userID, contentSHA).
		First(&letter).Error
	if err != nil {
		return nil, err
	}
	return &letter, nil
}

// Update updates an existing reference letter
func (r *GormReferenceLetterRepository) Update(ctx context.Context, letter *domain.ReferenceLetter) error {
	return r.db.WithContext(ctx).Save(letter).Error
}

// Delete deletes a reference letter
func (r *GormReferenceLetterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.ReferenceLetter{}, "id = ?", id).Error
}
