package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Skill represents a skill in the system
type Skill struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key"`
	Name string    `gorm:"uniqueIndex;not null"`
}

// BeforeCreate hook to generate UUID if not set
func (s *Skill) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// ProfileSkill represents the many-to-many relationship between Profile and Skill
type ProfileSkill struct {
	ProfileID   uuid.UUID `gorm:"type:uuid;primary_key"`
	SkillID     uuid.UUID `gorm:"type:uuid;primary_key"`
	Proficiency string    `gorm:"type:varchar(50)"`
}
