package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Profile represents a user's professional profile
type Profile struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Summary   string    `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	User            *User            `gorm:"foreignKey:UserID"`
	WorkExperiences []WorkExperience `gorm:"foreignKey:ProfileID"`
	Skills          []*Skill         `gorm:"many2many:profile_skills;"`
	JobMatches      []JobMatch       `gorm:"foreignKey:BaseProfileID"`
}

// Validate validates the profile data
func (p *Profile) Validate() error {
	if p.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	return nil
}

// NewProfile creates a new profile with a generated UUID
func NewProfile(userID uuid.UUID) *Profile {
	return &Profile{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AddWorkExperience adds a work experience to the profile
func (p *Profile) AddWorkExperience(we *WorkExperience) {
	we.ProfileID = p.ID
	p.WorkExperiences = append(p.WorkExperiences, *we)
}

// AddSkill adds a skill to the profile
func (p *Profile) AddSkill(skill *Skill) {
	p.Skills = append(p.Skills, skill)
}
