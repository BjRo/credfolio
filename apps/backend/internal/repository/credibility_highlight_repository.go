package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// CredibilityHighlightRepository defines the interface for credibility highlight data access
type CredibilityHighlightRepository interface {
	Create(ctx context.Context, highlight *domain.CredibilityHighlight) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CredibilityHighlight, error)
	GetByWorkExperienceID(ctx context.Context, workExpID uuid.UUID) ([]*domain.CredibilityHighlight, error)
	Update(ctx context.Context, highlight *domain.CredibilityHighlight) error
	Delete(ctx context.Context, id uuid.UUID) error
}
