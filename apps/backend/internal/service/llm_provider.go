package service

import (
	"context"
)

// ExtractedProfileData represents structured data extracted from a reference letter
type ExtractedProfileData struct {
	CompanyName  string   `json:"companyName" jsonschema_description:"The name of the company or organization"`
	Role         string   `json:"role" jsonschema_description:"The job title or role position"`
	StartDate    string   `json:"startDate" jsonschema_description:"Employment start date in YYYY-MM-DD format"`
	EndDate      string   `json:"endDate" jsonschema_description:"Employment end date in YYYY-MM-DD format, or empty string if current position"`
	Skills       []string `json:"skills" jsonschema_description:"List of skills mentioned or demonstrated"`
	Achievements []string `json:"achievements" jsonschema_description:"List of notable achievements or accomplishments"`
	Description  string   `json:"description" jsonschema_description:"Summary description of the role and responsibilities"`
}

// CredibilityData represents sentiment and quotes extracted from a reference letter
type CredibilityData struct {
	Quotes    []string `json:"quotes" jsonschema_description:"Positive quotes or praise from the reference letter"`
	Sentiment string   `json:"sentiment" jsonschema:"enum=POSITIVE,enum=NEUTRAL" jsonschema_description:"Overall sentiment of the reference letter"`
}

// LLMProvider defines the interface for LLM operations
type LLMProvider interface {
	// ExtractProfileData extracts structured profile data from text
	ExtractProfileData(ctx context.Context, text string) (*ExtractedProfileData, error)

	// ExtractCredibility extracts positive quotes and sentiment from text
	ExtractCredibility(ctx context.Context, text string) (*CredibilityData, error)

	// TailorProfile generates a tailored profile summary based on job description
	TailorProfile(ctx context.Context, profileText string, jobDescription string) (string, float64, error)
}
