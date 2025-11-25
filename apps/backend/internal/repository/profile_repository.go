package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// ProfileRepository defines the interface for profile data access
type ProfileRepository interface {
	// Create creates a new profile
	Create(ctx context.Context, profile *domain.Profile) error

	// Update updates an existing profile
	Update(ctx context.Context, profile *domain.Profile) error

	// FindByID finds a profile by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error)

	// FindByUserID finds a profile by user ID
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error)

	// Delete deletes a profile
	Delete(ctx context.Context, id uuid.UUID) error

	// FindByUserIDWithRelations finds a profile with all relations loaded
	FindByUserIDWithRelations(ctx context.Context, userID uuid.UUID) (*domain.Profile, error)
}
