package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// WorkExperienceRepository defines the interface for work experience data access
type WorkExperienceRepository interface {
	Create(ctx context.Context, exp *domain.WorkExperience) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.WorkExperience, error)
	GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.WorkExperience, error)
	Update(ctx context.Context, exp *domain.WorkExperience) error
	Delete(ctx context.Context, id uuid.UUID) error
}
