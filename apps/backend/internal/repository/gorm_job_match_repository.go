package repository

import (
	"context"
	"errors"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormJobMatchRepository implements JobMatchRepository using GORM
type GormJobMatchRepository struct {
	db *gorm.DB
}

// NewGormJobMatchRepository creates a new GORM-based job match repository
func NewGormJobMatchRepository(db *gorm.DB) *GormJobMatchRepository {
	return &GormJobMatchRepository{db: db}
}

// Create stores a new job match
func (r *GormJobMatchRepository) Create(ctx context.Context, jobMatch *domain.JobMatch) error {
	return r.db.WithContext(ctx).Create(jobMatch).Error
}

// FindByID retrieves a job match by its ID
func (r *GormJobMatchRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.JobMatch, error) {
	var jobMatch domain.JobMatch
	err := r.db.WithContext(ctx).
		Preload("BaseProfile").
		First(&jobMatch, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &jobMatch, nil
}

// FindByProfileID retrieves all job matches for a profile
func (r *GormJobMatchRepository) FindByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.JobMatch, error) {
	var jobMatches []*domain.JobMatch
	err := r.db.WithContext(ctx).
		Where("base_profile_id = ?", profileID).
		Order("created_at DESC").
		Find(&jobMatches).Error

	if err != nil {
		return nil, err
	}
	return jobMatches, nil
}

// Update updates an existing job match
func (r *GormJobMatchRepository) Update(ctx context.Context, jobMatch *domain.JobMatch) error {
	return r.db.WithContext(ctx).Save(jobMatch).Error
}

// Delete removes a job match
func (r *GormJobMatchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.JobMatch{}, "id = ?", id).Error
}
