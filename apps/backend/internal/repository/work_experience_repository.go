package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// WorkExperienceRepository defines the interface for work experience data access
type WorkExperienceRepository interface {
	// Create creates a new work experience
	Create(ctx context.Context, experience *domain.WorkExperience) error

	// Update updates an existing work experience
	Update(ctx context.Context, experience *domain.WorkExperience) error

	// FindByID finds a work experience by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*domain.WorkExperience, error)

	// FindByProfileID finds all work experiences for a profile
	FindByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.WorkExperience, error)

	// FindByProfileIDWithHighlights finds work experiences with credibility highlights
	FindByProfileIDWithHighlights(ctx context.Context, profileID uuid.UUID) ([]*domain.WorkExperience, error)

	// Delete deletes a work experience
	Delete(ctx context.Context, id uuid.UUID) error
}
