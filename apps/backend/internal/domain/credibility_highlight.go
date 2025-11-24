package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Sentiment represents the sentiment of a credibility highlight
type Sentiment string

const (
	SentimentPositive Sentiment = "POSITIVE"
	SentimentNeutral  Sentiment = "NEUTRAL"
)

// CredibilityHighlight represents a positive quote or sentiment from a reference letter
type CredibilityHighlight struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key"`
	WorkExperienceID uuid.UUID `gorm:"type:uuid;not null;index"`
	Quote            string    `gorm:"type:text;not null"`
	Sentiment        Sentiment `gorm:"type:varchar(20);not null;default:'POSITIVE'"`
	SourceLetterID   uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// BeforeCreate hook to generate UUID if not set
func (c *CredibilityHighlight) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
