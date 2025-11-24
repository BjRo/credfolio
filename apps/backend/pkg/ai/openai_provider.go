package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/credfolio/apps/backend/internal/service"
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIProvider implements the LLMProvider interface using OpenAI
type OpenAIProvider struct {
	client openai.Client
	model  string
}

// GenerateSchema generates a JSON schema for a given Go struct type
// Structured Outputs uses a subset of JSON schema
func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

// Generate schemas at initialization time
var (
	extractedProfileDataSchema = GenerateSchema[service.ExtractedProfileData]()
	credibilityDataSchema      = GenerateSchema[service.CredibilityData]()
)

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string, model string) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))

	// Use a model that supports structured outputs
	if model == "" {
		// gpt-4o-mini-2024-07-18 supports structured outputs and is cost-effective
		model = openai.ChatModelGPT4oMini2024_07_18
	}

	return &OpenAIProvider{
		client: client,
		model:  model,
	}, nil
}

// ExtractProfileData extracts structured profile data from text using OpenAI with structured outputs
func (p *OpenAIProvider) ExtractProfileData(ctx context.Context, text string) (*service.ExtractedProfileData, error) {
	systemMessage := `You are an expert data extraction assistant specializing in parsing professional reference letters and employment documents. Your role is to accurately extract structured information about work experience, including company names, job roles, employment dates, skills, achievements, and role descriptions.

Extract information precisely and completely. If a field cannot be determined from the text, use an empty string for text fields or an empty array for list fields. Dates must be in YYYY-MM-DD format. If the position is current, leave endDate as an empty string.`

	userMessage := fmt.Sprintf(`Extract structured information from the following reference letter text:

%s

Extract all available information about the work experience, including company name, role/title, start and end dates, skills mentioned, achievements, and a summary description of the role and responsibilities.`, text)

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "extracted_profile_data",
		Description: openai.String("Structured profile data extracted from a reference letter"),
		Schema:      extractedProfileDataSchema,
		Strict:      openai.Bool(true),
	}

	req := openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMessage),
			openai.UserMessage(userMessage),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
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
		return nil, fmt.Errorf("failed to parse OpenAI response: %w (content: %s)", err, content)
	}

	return &data, nil
}

// ExtractCredibility extracts positive quotes and sentiment from text using structured outputs
func (p *OpenAIProvider) ExtractCredibility(ctx context.Context, text string) (*service.CredibilityData, error) {
	systemMessage := `You are an expert sentiment analysis assistant specializing in analyzing professional reference letters. Your role is to identify positive feedback, praise, and credibility indicators from reference letters.

Extract specific positive quotes that highlight the person's strengths, achievements, or positive attributes. Determine the overall sentiment as either "POSITIVE" (if the letter contains clear positive feedback) or "NEUTRAL" (if the letter is factual without strong positive or negative sentiment).`

	userMessage := fmt.Sprintf(`Analyze the following reference letter text and extract positive quotes and sentiment:

%s

Identify specific positive quotes that demonstrate credibility, strengths, or achievements. Determine if the overall sentiment is POSITIVE or NEUTRAL.`, text)

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "credibility_data",
		Description: openai.String("Credibility indicators and sentiment extracted from a reference letter"),
		Schema:      credibilityDataSchema,
		Strict:      openai.Bool(true),
	}

	req := openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMessage),
			openai.UserMessage(userMessage),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
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
		return nil, fmt.Errorf("failed to parse OpenAI response: %w (content: %s)", err, content)
	}

	return &data, nil
}

// TailorProfile generates a tailored profile summary based on job description using structured outputs
func (p *OpenAIProvider) TailorProfile(ctx context.Context, profileText string, jobDescription string) (string, float64, error) {
	systemMessage := `You are an expert career advisor and resume tailoring specialist. Your role is to analyze professional profiles and job descriptions to create tailored summaries that highlight the most relevant experience and skills.

Generate a compelling summary that emphasizes how the candidate's experience aligns with the job requirements. Provide a match score from 0.0 to 1.0 that reflects how well the profile matches the job description, where 1.0 represents a perfect match.`

	userMessage := fmt.Sprintf(`Given this professional profile:

%s

And this job description:

%s

Generate a tailored summary that emphasizes relevant experience and skills, and provide a match score from 0.0 to 1.0.`, profileText, jobDescription)

	type TailorResult struct {
		Summary    string  `json:"summary" jsonschema_description:"Tailored profile summary emphasizing relevant experience"`
		MatchScore float64 `json:"matchScore" jsonschema_description:"Match score from 0.0 to 1.0 indicating how well the profile matches the job"`
	}

	tailorSchema := GenerateSchema[TailorResult]()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "tailor_result",
		Description: openai.String("Tailored profile summary and match score"),
		Schema:      tailorSchema,
		Strict:      openai.Bool(true),
	}

	req := openai.ChatCompletionNewParams{
		Model: p.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMessage),
			openai.UserMessage(userMessage),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
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
	var result TailorResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return "", 0, fmt.Errorf("failed to parse OpenAI response: %w (content: %s)", err, content)
	}

	return result.Summary, result.MatchScore, nil
}
