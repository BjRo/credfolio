package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/credfolio/apps/backend/internal/service"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIProvider implements the LLMProvider interface using OpenAI
type OpenAIProvider struct {
	client openai.Client
	model  string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string, model string) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	if model == "" {
		model = openai.ChatModelGPT4oMini // Default to cost-effective model
	}

	return &OpenAIProvider{
		client: client,
		model:  model,
	}, nil
}

// ExtractProfileData extracts structured profile data from text using OpenAI
func (p *OpenAIProvider) ExtractProfileData(ctx context.Context, text string) (*service.ExtractedProfileData, error) {
	prompt := fmt.Sprintf(`Extract the following information from this reference letter text in JSON format:
{
  "companyName": "Company name",
  "role": "Job title/role",
  "startDate": "YYYY-MM-DD",
  "endDate": "YYYY-MM-DD or empty string if current",
  "skills": ["skill1", "skill2"],
  "achievements": ["achievement1", "achievement2"],
  "description": "Job description summary"
}

Reference letter text:
%s

Return only valid JSON, no other text.`, text)

	req := openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	}

	resp, err := p.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content
	var data service.ExtractedProfileData
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	return &data, nil
}

// ExtractCredibility extracts positive quotes and sentiment from text
func (p *OpenAIProvider) ExtractCredibility(ctx context.Context, text string) (*service.CredibilityData, error) {
	prompt := fmt.Sprintf(`Extract positive quotes and sentiment from this reference letter text. Return JSON:
{
  "quotes": ["quote1", "quote2"],
  "sentiment": "POSITIVE" or "NEUTRAL"
}

Reference letter text:
%s

Return only valid JSON, no other text.`, text)

	req := openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	}

	resp, err := p.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content
	var data service.CredibilityData
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	return &data, nil
}

// TailorProfile generates a tailored profile summary based on job description
func (p *OpenAIProvider) TailorProfile(ctx context.Context, profileText string, jobDescription string) (string, float64, error) {
	prompt := fmt.Sprintf(`Given this profile:
%s

And this job description:
%s

Generate a tailored summary that emphasizes relevant experience and skills. Also provide a match score from 0.0 to 1.0.

Return JSON:
{
  "summary": "tailored summary text",
  "matchScore": 0.85
}

Return only valid JSON, no other text.`, profileText, jobDescription)

	req := openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	}

	resp, err := p.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to call OpenAI: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", 0, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content
	var result struct {
		Summary    string  `json:"summary"`
		MatchScore float64 `json:"matchScore"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return "", 0, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	return result.Summary, result.MatchScore, nil
}
