package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// JobMatchRepository defines the interface for job match data access
type JobMatchRepository interface {
	Create(ctx context.Context, jobMatch *domain.JobMatch) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.JobMatch, error)
	GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.JobMatch, error)
	Update(ctx context.Context, jobMatch *domain.JobMatch) error
	Delete(ctx context.Context, id uuid.UUID) error
}
