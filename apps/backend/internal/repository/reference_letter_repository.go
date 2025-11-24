package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// ReferenceLetterRepository defines the interface for reference letter data access
type ReferenceLetterRepository interface {
	Create(ctx context.Context, letter *domain.ReferenceLetter) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ReferenceLetter, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error)
	Update(ctx context.Context, letter *domain.ReferenceLetter) error
	Delete(ctx context.Context, id uuid.UUID) error
}
