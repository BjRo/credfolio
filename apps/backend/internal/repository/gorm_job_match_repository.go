package repository

import (
	"context"

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

// Create creates a new job match
func (r *GormJobMatchRepository) Create(ctx context.Context, jobMatch *domain.JobMatch) error {
	return r.db.WithContext(ctx).Create(jobMatch).Error
}

// GetByID retrieves a job match by ID
func (r *GormJobMatchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.JobMatch, error) {
	var jobMatch domain.JobMatch
	err := r.db.WithContext(ctx).First(&jobMatch, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &jobMatch, nil
}

// GetByProfileID retrieves all job matches for a profile
func (r *GormJobMatchRepository) GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.JobMatch, error) {
	var jobMatches []*domain.JobMatch
	err := r.db.WithContext(ctx).Where("base_profile_id = ?", profileID).Find(&jobMatches).Error
	if err != nil {
		return nil, err
	}
	return jobMatches, nil
}

// Update updates an existing job match
func (r *GormJobMatchRepository) Update(ctx context.Context, jobMatch *domain.JobMatch) error {
	return r.db.WithContext(ctx).Save(jobMatch).Error
}

// Delete deletes a job match
func (r *GormJobMatchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.JobMatch{}, "id = ?", id).Error
}
