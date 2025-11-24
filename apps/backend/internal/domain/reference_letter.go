package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReferenceLetterStatus represents the processing status of a reference letter
type ReferenceLetterStatus string

const (
	ReferenceLetterStatusPending   ReferenceLetterStatus = "PENDING"
	ReferenceLetterStatusProcessed ReferenceLetterStatus = "PROCESSED"
	ReferenceLetterStatusFailed    ReferenceLetterStatus = "FAILED"
)

// ReferenceLetter represents an uploaded reference letter
type ReferenceLetter struct {
	ID            uuid.UUID             `gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID             `gorm:"type:uuid;not null;index"`
	FileName      string                `gorm:"not null"`
	StoragePath   string                `gorm:"not null"`
	UploadDate    time.Time             `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Status        ReferenceLetterStatus `gorm:"type:varchar(20);not null;default:'PENDING'"`
	ExtractedText string                `gorm:"type:text"`
	ContentSHA    string                `gorm:"type:varchar(64);index"` // SHA256 hash of ExtractedText for deduplication
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// BeforeCreate hook to generate UUID if not set
func (r *ReferenceLetter) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
