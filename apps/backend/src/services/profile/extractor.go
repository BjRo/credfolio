package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/credfolio/apps/backend/src/services/llm"
)

// Interfaces for dependencies
type PDFExtractor interface {
	ExtractText(r io.ReaderAt, size int64) (string, error)
}

type LLMClient interface {
	GenerateJSON(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

type Extractor struct {
	pdfExtractor PDFExtractor
	llmClient    LLMClient
}

func NewExtractor(pdf PDFExtractor, llm LLMClient) *Extractor {
	return &Extractor{
		pdfExtractor: pdf,
		llmClient:    llm,
	}
}

// ExtractedData matches the JSON schema from the prompt
type ExtractedData struct {
	FullName string `json:"full_name"`
	Company  struct {
		Name      string  `json:"name"`
		StartDate string  `json:"start_date"`
		EndDate   *string `json:"end_date"`
		LogoURL   *string `json:"logo_url"`
	} `json:"company"`
	Role struct {
		Title            string   `json:"title"`
		Description      string   `json:"description"`
		EmployerFeedback string   `json:"employer_feedback"`
		Skills           []string `json:"skills"`
	} `json:"role"`
}

func (s *Extractor) Extract(ctx context.Context, file io.ReaderAt, size int64) (*ExtractedData, error) {
	text, err := s.pdfExtractor.ExtractText(file, size)
	if err != nil {
		return nil, fmt.Errorf("pdf extraction failed: %w", err)
	}

	// Prompt is imported from llm package
	jsonStr, err := s.llmClient.GenerateJSON(ctx, llm.ReferenceExtractionSystemPrompt, text)
	if err != nil {
		return nil, fmt.Errorf("llm extraction failed: %w", err)
	}

	var data ExtractedData
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to parse llm response: %w", err)
	}

	return &data, nil
}
