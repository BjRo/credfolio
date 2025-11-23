package profile

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/credfolio/apps/backend/src/services/llm"
	"github.com/google/uuid"
)

type TailoringResult struct {
	RelevantSkillIDs      []string `json:"relevant_skill_ids"`
	RelevantExperienceIDs []string `json:"relevant_experience_ids"`
	SummaryHighlights     string   `json:"summary_highlights"`
	MatchScore            int      `json:"match_score"`
}

func (s *Service) TailorProfile(ctx context.Context, userID uuid.UUID, jobDescription string) (*TailoringResult, error) {
	profile, err := s.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal profile: %w", err)
	}

	userPrompt := fmt.Sprintf("Profile JSON:\n%s\n\nJob Description:\n%s", string(profileJSON), jobDescription)

	jsonStr, err := s.llmClient.GenerateJSON(ctx, llm.TailoringSystemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("llm tailoring failed: %w", err)
	}

	var result TailoringResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse tailoring result: %w", err)
	}

	return &result, nil
}
