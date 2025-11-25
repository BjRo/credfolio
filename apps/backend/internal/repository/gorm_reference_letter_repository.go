package repository

import (
	"context"
	"errors"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormReferenceLetterRepository implements ReferenceLetterRepository using GORM
type GormReferenceLetterRepository struct {
	db *gorm.DB
}

// NewGormReferenceLetterRepository creates a new GORM reference letter repository
func NewGormReferenceLetterRepository(db *gorm.DB) *GormReferenceLetterRepository {
	return &GormReferenceLetterRepository{db: db}
}

// Create creates a new reference letter
func (r *GormReferenceLetterRepository) Create(ctx context.Context, letter *domain.ReferenceLetter) error {
	return r.db.WithContext(ctx).Create(letter).Error
}

// Update updates an existing reference letter
func (r *GormReferenceLetterRepository) Update(ctx context.Context, letter *domain.ReferenceLetter) error {
	return r.db.WithContext(ctx).Save(letter).Error
}

// FindByID finds a reference letter by its ID
func (r *GormReferenceLetterRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ReferenceLetter, error) {
	var letter domain.ReferenceLetter
	err := r.db.WithContext(ctx).First(&letter, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &letter, nil
}

// FindByUserID finds all reference letters for a user
func (r *GormReferenceLetterRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	var letters []*domain.ReferenceLetter
	err := r.db.WithContext(ctx).Find(&letters, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return letters, nil
}

// FindPendingByUserID finds pending reference letters for a user
func (r *GormReferenceLetterRepository) FindPendingByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	var letters []*domain.ReferenceLetter
	err := r.db.WithContext(ctx).Find(&letters, "user_id = ? AND status = ?", userID, domain.StatusPending).Error
	if err != nil {
		return nil, err
	}
	return letters, nil
}

// Delete deletes a reference letter
func (r *GormReferenceLetterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.ReferenceLetter{}, "id = ?", id).Error
}

// Ensure GormReferenceLetterRepository implements ReferenceLetterRepository
var _ ReferenceLetterRepository = (*GormReferenceLetterRepository)(nil)
