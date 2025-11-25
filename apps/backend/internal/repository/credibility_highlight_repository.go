package repository

import (
	"context"

	"github.com/credfolio/apps/backend/internal/domain"
	"github.com/google/uuid"
)

// CredibilityHighlightRepository defines the interface for credibility highlight data access
type CredibilityHighlightRepository interface {
	// Create creates a new credibility highlight
	Create(ctx context.Context, highlight *domain.CredibilityHighlight) error

	// Update updates an existing credibility highlight
	Update(ctx context.Context, highlight *domain.CredibilityHighlight) error

	// FindByID finds a credibility highlight by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*domain.CredibilityHighlight, error)

	// FindByWorkExperienceID finds all highlights for a work experience
	FindByWorkExperienceID(ctx context.Context, workExperienceID uuid.UUID) ([]*domain.CredibilityHighlight, error)

	// FindBySourceLetterID finds all highlights from a reference letter
	FindBySourceLetterID(ctx context.Context, sourceLetterID uuid.UUID) ([]*domain.CredibilityHighlight, error)

	// Delete deletes a credibility highlight
	Delete(ctx context.Context, id uuid.UUID) error
}
