package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JobMatch represents a tailored profile for a specific job description
type JobMatch struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	BaseProfileID   uuid.UUID `gorm:"type:uuid;not null;index"`
	JobDescription  string    `gorm:"type:text;not null"`
	MatchScore      float64   `gorm:"not null"`
	TailoredSummary string    `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// BeforeCreate hook to generate UUID if not set
func (j *JobMatch) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
