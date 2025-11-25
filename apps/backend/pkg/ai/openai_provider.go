package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/credfolio/apps/backend/internal/service"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// OpenAIProvider implements LLMProvider using OpenAI API
type OpenAIProvider struct {
	client openai.Client
	model  openai.ChatModel
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIProvider{
		client: client,
		model:  openai.ChatModelGPT4oMini,
	}
}

// ExtractProfileData extracts structured profile data from reference letter text
func (p *OpenAIProvider) ExtractProfileData(ctx context.Context, referenceText string) (*service.ProfileData, error) {
	if referenceText == "" {
		return nil, errors.New("reference text is required")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	prompt := fmt.Sprintf(`Extract structured profile data from the following reference letter. Return a JSON object with:
- summary: A brief professional summary
- work_experiences: Array of {company_name, role, description}
- skills: Array of skill names mentioned
- highlights: Array of {quote, sentiment} where sentiment is "POSITIVE" or "NEUTRAL"

Reference letter:
%s`, referenceText)

	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a professional resume parser. Extract structured data from reference letters and return valid JSON."),
			openai.UserMessage(prompt),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openai api error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	var result service.ProfileData
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// TailorProfile tailors a profile to a job description
func (p *OpenAIProvider) TailorProfile(ctx context.Context, profileSummary string, experiences []service.WorkExperience, skills []string, jobDescription string) (*service.TailoredProfileData, error) {
	if jobDescription == "" {
		return nil, errors.New("job description is required")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	experiencesJSON, _ := json.Marshal(experiences)
	skillsJSON, _ := json.Marshal(skills)

	prompt := fmt.Sprintf(`Given this profile and job description, create a tailored profile that emphasizes relevant experience. Return a JSON object with:
- tailored_summary: A summary tailored to the job
- match_score: A float between 0 and 1 indicating how well the profile matches
- highlighted_skills: Array of skills most relevant to the job
- match_explanation: Brief explanation of why this is a good/poor match

Profile Summary: %s
Work Experiences: %s
Skills: %s

Job Description:
%s`, profileSummary, experiencesJSON, skillsJSON, jobDescription)

	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a career advisor helping tailor resumes to job descriptions. Return valid JSON."),
			openai.UserMessage(prompt),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openai api error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	var result service.TailoredProfileData
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// ExtractCredibilityHighlights extracts positive sentiment quotes from reference text
func (p *OpenAIProvider) ExtractCredibilityHighlights(ctx context.Context, referenceText string) ([]service.HighlightData, error) {
	if referenceText == "" {
		return nil, errors.New("reference text is required")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	prompt := fmt.Sprintf(`Extract positive quotes and testimonials from this reference letter. Return a JSON object with:
- highlights: Array of {quote, sentiment} where sentiment is "POSITIVE" or "NEUTRAL"

Focus on quotes that demonstrate:
- Professional achievements
- Personal qualities
- Team collaboration
- Technical skills
- Leadership

Reference letter:
%s`, referenceText)

	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are extracting credibility highlights from reference letters. Return valid JSON with a 'highlights' array."),
			openai.UserMessage(prompt),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openai api error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	var result struct {
		Highlights []service.HighlightData `json:"highlights"`
	}
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Highlights, nil
}

// CalculateRelevance calculates the relevance score between text and a reference
func (p *OpenAIProvider) CalculateRelevance(ctx context.Context, text string, reference string) (float64, error) {
	if text == "" || reference == "" {
		return 0, errors.New("both text and reference are required")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	prompt := fmt.Sprintf(`Rate the relevance of this text to the job description on a scale of 0 to 1.
Return a JSON object with a single field "score" containing a float between 0 and 1.

Text to evaluate:
%s

Job Description:
%s`, text, reference)

	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are evaluating how relevant work experience is to a job description. Return valid JSON with a 'score' field."),
			openai.UserMessage(prompt),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
		},
	})
	if err != nil {
		return 0, fmt.Errorf("openai api error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return 0, errors.New("no response from OpenAI")
	}

	var result struct {
		Score float64 `json:"score"`
	}
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	// Ensure score is between 0 and 1
	if result.Score < 0 {
		result.Score = 0
	}
	if result.Score > 1 {
		result.Score = 1
	}

	return result.Score, nil
}

// GenerateText generates text based on a prompt
func (p *OpenAIProvider) GenerateText(ctx context.Context, prompt string) (string, error) {
	if prompt == "" {
		return "", errors.New("prompt is required")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful assistant that generates clear, concise text."),
			openai.UserMessage(prompt),
		},
	})
	if err != nil {
		return "", fmt.Errorf("openai api error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// Ensure OpenAIProvider implements LLMProvider
var _ service.LLMProvider = (*OpenAIProvider)(nil)
