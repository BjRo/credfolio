package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// JobMatch represents a tailored profile for a specific job description
type JobMatch struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	BaseProfileID   uuid.UUID `gorm:"type:uuid;not null;index"`
	JobDescription  string    `gorm:"type:text;not null"`
	MatchScore      float64   `gorm:"type:decimal(3,2);not null"`
	TailoredSummary string    `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Relationships
	BaseProfile *Profile `gorm:"foreignKey:BaseProfileID"`
}

// Validate validates the job match data
func (jm *JobMatch) Validate() error {
	if jm.BaseProfileID == uuid.Nil {
		return errors.New("base profile ID is required")
	}
	if jm.JobDescription == "" {
		return errors.New("job description is required")
	}
	if jm.MatchScore < 0 || jm.MatchScore > 1 {
		return errors.New("match score must be between 0 and 1")
	}
	return nil
}

// NewJobMatch creates a new job match with a generated UUID
func NewJobMatch(baseProfileID uuid.UUID, jobDescription string) *JobMatch {
	return &JobMatch{
		ID:             uuid.New(),
		BaseProfileID:  baseProfileID,
		JobDescription: jobDescription,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// SetMatchResult sets the match score and tailored summary
func (jm *JobMatch) SetMatchResult(score float64, tailoredSummary string) {
	jm.MatchScore = score
	jm.TailoredSummary = tailoredSummary
	jm.UpdatedAt = time.Now()
}

// IsHighMatch returns true if the match score is above 0.7
func (jm *JobMatch) IsHighMatch() bool {
	return jm.MatchScore >= 0.7
}

// IsMediumMatch returns true if the match score is between 0.4 and 0.7
func (jm *JobMatch) IsMediumMatch() bool {
	return jm.MatchScore >= 0.4 && jm.MatchScore < 0.7
}

// IsLowMatch returns true if the match score is below 0.4
func (jm *JobMatch) IsLowMatch() bool {
	return jm.MatchScore < 0.4
}
