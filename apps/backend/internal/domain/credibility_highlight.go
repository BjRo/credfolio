package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Sentiment represents the sentiment of a credibility highlight
type Sentiment string

const (
	SentimentPositive Sentiment = "POSITIVE"
	SentimentNeutral  Sentiment = "NEUTRAL"
)

// CredibilityHighlight represents a quote or highlight from a reference letter
type CredibilityHighlight struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key"`
	WorkExperienceID uuid.UUID `gorm:"type:uuid;not null;index"`
	Quote            string    `gorm:"type:text;not null"`
	Sentiment        Sentiment `gorm:"type:varchar(20);not null"`
	SourceLetterID   uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Relationships
	WorkExperience *WorkExperience  `gorm:"foreignKey:WorkExperienceID"`
	SourceLetter   *ReferenceLetter `gorm:"foreignKey:SourceLetterID"`
}

// Validate validates the credibility highlight data
func (ch *CredibilityHighlight) Validate() error {
	if ch.Quote == "" {
		return errors.New("quote is required")
	}
	if ch.Sentiment != SentimentPositive && ch.Sentiment != SentimentNeutral {
		return errors.New("sentiment must be POSITIVE or NEUTRAL")
	}
	if ch.SourceLetterID == uuid.Nil {
		return errors.New("source letter ID is required")
	}
	return nil
}

// NewCredibilityHighlight creates a new credibility highlight with a generated UUID
func NewCredibilityHighlight(quote string, sentiment Sentiment, sourceLetterID uuid.UUID) *CredibilityHighlight {
	return &CredibilityHighlight{
		ID:             uuid.New(),
		Quote:          quote,
		Sentiment:      sentiment,
		SourceLetterID: sourceLetterID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// IsPositive returns true if the sentiment is positive
func (ch *CredibilityHighlight) IsPositive() bool {
	return ch.Sentiment == SentimentPositive
}
