package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Profile represents a user's professional profile
type Profile struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null;index"`
	Summary         string           `gorm:"type:text"`
	WorkExperiences []WorkExperience `gorm:"foreignKey:ProfileID;constraint:OnDelete:CASCADE"`
	Skills          []*Skill         `gorm:"many2many:profile_skills;"`
	JobMatches      []JobMatch       `gorm:"foreignKey:BaseProfileID;constraint:OnDelete:CASCADE"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// BeforeCreate hook to generate UUID if not set
func (p *Profile) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
