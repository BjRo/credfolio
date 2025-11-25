package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Skill represents a skill that can be associated with profiles
type Skill struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Name      string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	Profiles []*Profile `gorm:"many2many:profile_skills;"`
}

// ProfileSkill represents the join table for Profile and Skill many-to-many relationship
type ProfileSkill struct {
	ProfileID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	SkillID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Proficiency string    `gorm:"type:varchar(50)"`
	CreatedAt   time.Time
}

// Validate validates the skill data
func (s *Skill) Validate() error {
	if s.Name == "" {
		return errors.New("skill name is required")
	}
	return nil
}

// NewSkill creates a new skill with a generated UUID
func NewSkill(name string) *Skill {
	return &Skill{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
