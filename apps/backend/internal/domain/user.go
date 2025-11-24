package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID    uuid.UUID `gorm:"type:uuid;primary_key"`
	Email string    `gorm:"uniqueIndex;not null"`
	Name  string    `gorm:"not null"`
}

// BeforeCreate hook to generate UUID if not set
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
