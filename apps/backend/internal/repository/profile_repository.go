package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// ProfileRepository defines the interface for profile data access
type ProfileRepository interface {
	Create(ctx context.Context, profile *domain.Profile) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error)
	Update(ctx context.Context, profile *domain.Profile) error
	Delete(ctx context.Context, id uuid.UUID) error
}
