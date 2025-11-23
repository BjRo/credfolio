package models

import (
	"time"

	"github.com/google/uuid"
)

// UserProfile corresponds to the 'users' table
type UserProfile struct {
	ID        uuid.UUID      `json:"id"`
	Email     string         `json:"email"`
	FullName  string         `json:"full_name"`
	CreatedAt time.Time      `json:"created_at"`
	Companies []CompanyEntry `json:"companies,omitempty"` // Helper for API responses
}

type ExtractionStatus string

const (
	StatusPending   ExtractionStatus = "PENDING"
	StatusCompleted ExtractionStatus = "COMPLETED"
	StatusFailed    ExtractionStatus = "FAILED"
)

// ReferenceLetter corresponds to the 'reference_letters' table
type ReferenceLetter struct {
	ID                uuid.UUID              `json:"id"`
	UserID            uuid.UUID              `json:"user_id"`
	Filename          string                 `json:"filename"`
	StoragePath       string                 `json:"storage_path"`
	ContentHash       string                 `json:"content_hash"`
	UploadDate        time.Time              `json:"upload_date"`
	ExtractionStatus  ExtractionStatus       `json:"extraction_status"`
	ExtractedMetadata map[string]interface{} `json:"extracted_metadata"` // JSONB
}

// CompanyEntry corresponds to the 'companies' table
type CompanyEntry struct {
	ID        uuid.UUID        `json:"id"`
	UserID    uuid.UUID        `json:"user_id"`
	Name      string           `json:"name"`
	LogoURL   *string          `json:"logo_url,omitempty"`
	StartDate time.Time        `json:"start_date"`
	EndDate   *time.Time       `json:"end_date,omitempty"`
	Roles     []WorkExperience `json:"roles,omitempty"` // Helper for API responses
}

type SourceType string

const (
	SourceVerified     SourceType = "VERIFIED"
	SourceSelfReported SourceType = "SELF_REPORTED"
)

// WorkExperience corresponds to the 'work_experiences' table
type WorkExperience struct {
	ID                uuid.UUID  `json:"id"`
	CompanyID         uuid.UUID  `json:"company_id"`
	Title             string     `json:"title"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           *time.Time `json:"end_date,omitempty"`
	Description       string     `json:"description"`
	Source            SourceType `json:"source"`
	EmployerFeedback  string     `json:"employer_feedback"`
	ReferenceLetterID *uuid.UUID `json:"reference_letter_id,omitempty"`
	IsVerified        bool       `json:"is_verified"`
	Skills            []Skill    `json:"skills,omitempty"` // Helper for API responses
}

// Skill corresponds to the 'skills' table
type Skill struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
