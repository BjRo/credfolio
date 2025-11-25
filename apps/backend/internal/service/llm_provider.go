package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ProfileData represents structured data extracted from a reference letter
type ProfileData struct {
	Summary         string           `json:"summary"`
	WorkExperiences []WorkExperience `json:"work_experiences"`
	Skills          []string         `json:"skills"`
	Highlights      []HighlightData  `json:"highlights"`
}

// WorkExperience represents extracted work experience data
type WorkExperience struct {
	CompanyName string     `json:"company_name"`
	Role        string     `json:"role"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Description string     `json:"description"`
}

// HighlightData represents extracted credibility highlight data
type HighlightData struct {
	Quote     string `json:"quote"`
	Sentiment string `json:"sentiment"`
}

// TailoredProfileData represents a profile tailored to a job description
type TailoredProfileData struct {
	TailoredSummary   string             `json:"tailored_summary"`
	MatchScore        float64            `json:"match_score"`
	RankedExperiences []RankedExperience `json:"ranked_experiences"`
	HighlightedSkills []string           `json:"highlighted_skills"`
	MatchExplanation  string             `json:"match_explanation"`
}

// RankedExperience represents a work experience ranked by relevance to job
type RankedExperience struct {
	ExperienceID   uuid.UUID `json:"experience_id"`
	RelevanceScore float64   `json:"relevance_score"`
	Explanation    string    `json:"explanation"`
}

// LLMProvider defines the interface for interacting with LLM services
type LLMProvider interface {
	// ExtractProfileData extracts structured profile data from reference letter text
	ExtractProfileData(ctx context.Context, referenceText string) (*ProfileData, error)

	// TailorProfile tailors a profile summary and experiences to a job description
	TailorProfile(ctx context.Context, profileSummary string, experiences []WorkExperience, skills []string, jobDescription string) (*TailoredProfileData, error)

	// ExtractCredibilityHighlights extracts positive sentiment quotes from reference text
	ExtractCredibilityHighlights(ctx context.Context, referenceText string) ([]HighlightData, error)

	// CalculateRelevance calculates the relevance score between text and a reference
	CalculateRelevance(ctx context.Context, text string, reference string) (float64, error)

	// GenerateText generates text based on a prompt
	GenerateText(ctx context.Context, prompt string) (string, error)
}
