package service

import (
	"context"
)

// ExtractedProfileData represents structured data extracted from a reference letter
type ExtractedProfileData struct {
	CompanyName  string
	Role         string
	StartDate    string
	EndDate      string
	Skills       []string
	Achievements []string
	Description  string
}

// CredibilityData represents sentiment and quotes extracted from a reference letter
type CredibilityData struct {
	Quotes    []string
	Sentiment string // "POSITIVE" or "NEUTRAL"
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
