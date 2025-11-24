package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WorkExperience represents a work experience entry
type WorkExperience struct {
	ID                    uuid.UUID              `gorm:"type:uuid;primary_key"`
	ProfileID             uuid.UUID              `gorm:"type:uuid;not null;index"`
	CompanyName           string                 `gorm:"not null"`
	Role                  string                 `gorm:"not null"`
	StartDate             time.Time              `gorm:"not null"`
	EndDate               *time.Time             `gorm:"index"`
	Description           string                 `gorm:"type:text"`
	ReferenceLetterID     *uuid.UUID             `gorm:"type:uuid;index"`
	CredibilityHighlights []CredibilityHighlight `gorm:"foreignKey:WorkExperienceID;constraint:OnDelete:CASCADE"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// BeforeCreate hook to generate UUID if not set
func (w *WorkExperience) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
