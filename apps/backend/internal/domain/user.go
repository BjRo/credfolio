package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system (mock/existing entity)
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	Profile          *Profile          `gorm:"foreignKey:UserID"`
	ReferenceLetters []ReferenceLetter `gorm:"foreignKey:UserID"`
}

// Validate validates the user data
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// NewUser creates a new user with a generated UUID
func NewUser(email, name string) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
