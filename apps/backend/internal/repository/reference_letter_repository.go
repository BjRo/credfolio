package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// ReferenceLetterRepository defines the interface for reference letter data access
type ReferenceLetterRepository interface {
	// Create creates a new reference letter
	Create(ctx context.Context, letter *domain.ReferenceLetter) error

	// Update updates an existing reference letter
	Update(ctx context.Context, letter *domain.ReferenceLetter) error

	// FindByID finds a reference letter by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*domain.ReferenceLetter, error)

	// FindByUserID finds all reference letters for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error)

	// FindPendingByUserID finds pending reference letters for a user
	FindPendingByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error)

	// Delete deletes a reference letter
	Delete(ctx context.Context, id uuid.UUID) error
}
