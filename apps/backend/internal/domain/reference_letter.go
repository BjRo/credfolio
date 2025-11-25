package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ReferenceLetterStatus represents the processing status of a reference letter
type ReferenceLetterStatus string

const (
	StatusPending   ReferenceLetterStatus = "PENDING"
	StatusProcessed ReferenceLetterStatus = "PROCESSED"
	StatusFailed    ReferenceLetterStatus = "FAILED"
)

// ReferenceLetter represents an uploaded reference letter
type ReferenceLetter struct {
	ID            uuid.UUID             `gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID             `gorm:"type:uuid;not null;index"`
	FileName      string                `gorm:"type:varchar(255);not null"`
	StoragePath   string                `gorm:"type:varchar(500);not null"`
	UploadDate    time.Time             `gorm:"not null"`
	Status        ReferenceLetterStatus `gorm:"type:varchar(20);not null;default:'PENDING'"`
	ExtractedText string                `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relationships
	User            *User            `gorm:"foreignKey:UserID"`
	WorkExperiences []WorkExperience `gorm:"foreignKey:ReferenceLetterID"`
}

// Validate validates the reference letter data
func (rl *ReferenceLetter) Validate() error {
	if rl.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if rl.FileName == "" {
		return errors.New("file name is required")
	}
	if rl.StoragePath == "" {
		return errors.New("storage path is required")
	}
	return nil
}

// NewReferenceLetter creates a new reference letter with a generated UUID
func NewReferenceLetter(userID uuid.UUID, fileName, storagePath string) *ReferenceLetter {
	return &ReferenceLetter{
		ID:          uuid.New(),
		UserID:      userID,
		FileName:    fileName,
		StoragePath: storagePath,
		UploadDate:  time.Now(),
		Status:      StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// MarkProcessed marks the reference letter as successfully processed
func (rl *ReferenceLetter) MarkProcessed(extractedText string) {
	rl.Status = StatusProcessed
	rl.ExtractedText = extractedText
	rl.UpdatedAt = time.Now()
}

// MarkFailed marks the reference letter as failed processing
func (rl *ReferenceLetter) MarkFailed() {
	rl.Status = StatusFailed
	rl.UpdatedAt = time.Now()
}

// IsProcessed returns true if the reference letter has been processed
func (rl *ReferenceLetter) IsProcessed() bool {
	return rl.Status == StatusProcessed
}

// IsPending returns true if the reference letter is pending processing
func (rl *ReferenceLetter) IsPending() bool {
	return rl.Status == StatusPending
}
