package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// JobMatchRepository defines the interface for job match data access
type JobMatchRepository interface {
	// Create stores a new job match
	Create(ctx context.Context, jobMatch *domain.JobMatch) error

	// FindByID retrieves a job match by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*domain.JobMatch, error)

	// FindByProfileID retrieves all job matches for a profile
	FindByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.JobMatch, error)

	// Update updates an existing job match
	Update(ctx context.Context, jobMatch *domain.JobMatch) error

	// Delete removes a job match
	Delete(ctx context.Context, id uuid.UUID) error
}
