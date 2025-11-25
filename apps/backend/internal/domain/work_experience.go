package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// WorkExperience represents a work experience entry in a profile
type WorkExperience struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key"`
	ProfileID         uuid.UUID  `gorm:"type:uuid;not null;index"`
	CompanyName       string     `gorm:"type:varchar(255);not null"`
	Role              string     `gorm:"type:varchar(255);not null"`
	StartDate         time.Time  `gorm:"type:date;not null"`
	EndDate           *time.Time `gorm:"type:date"`
	Description       string     `gorm:"type:text"`
	ReferenceLetterID *uuid.UUID `gorm:"type:uuid"`
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Relationships
	Profile               *Profile               `gorm:"foreignKey:ProfileID"`
	ReferenceLetter       *ReferenceLetter       `gorm:"foreignKey:ReferenceLetterID"`
	CredibilityHighlights []CredibilityHighlight `gorm:"foreignKey:WorkExperienceID"`
}

// Validate validates the work experience data
func (we *WorkExperience) Validate() error {
	if we.CompanyName == "" {
		return errors.New("company name is required")
	}
	if we.Role == "" {
		return errors.New("role is required")
	}
	if we.StartDate.IsZero() {
		return errors.New("start date is required")
	}
	if we.EndDate != nil && we.EndDate.Before(we.StartDate) {
		return errors.New("end date must be after start date")
	}
	return nil
}

// NewWorkExperience creates a new work experience with a generated UUID
func NewWorkExperience(companyName, role string) *WorkExperience {
	return &WorkExperience{
		ID:          uuid.New(),
		CompanyName: companyName,
		Role:        role,
		StartDate:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// SetDates sets the start and end dates for the work experience
func (we *WorkExperience) SetDates(startDate time.Time, endDate *time.Time) {
	we.StartDate = startDate
	we.EndDate = endDate
}

// IsCurrent returns true if this is a current position (no end date)
func (we *WorkExperience) IsCurrent() bool {
	return we.EndDate == nil
}

// AddCredibilityHighlight adds a credibility highlight to the work experience
func (we *WorkExperience) AddCredibilityHighlight(ch *CredibilityHighlight) {
	ch.WorkExperienceID = we.ID
	we.CredibilityHighlights = append(we.CredibilityHighlights, *ch)
}
