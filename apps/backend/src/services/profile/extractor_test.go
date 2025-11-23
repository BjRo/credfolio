package profile

import (
	"context"
	"io"
	"strings"
	"testing"
)

type MockPDFExtractor struct {
	Text string
	Err  error
}

func (m *MockPDFExtractor) ExtractText(r io.ReaderAt, size int64) (string, error) {
	return m.Text, m.Err
}

type MockLLMClient struct {
	Response string
	Err      error
}

func (m *MockLLMClient) GenerateJSON(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	return m.Response, m.Err
}

func TestExtract(t *testing.T) {
	pdfText := "John Doe worked at Acme Corp from 2020 to 2022."
	jsonResp := `{
		"full_name": "John Doe",
		"company": {"name": "Acme Corp", "start_date": "2020-01-01", "end_date": "2022-01-01"},
		"role": {"title": "Engineer", "description": "Dev", "employer_feedback": "Good", "skills": ["Go"]}
	}`

	mockPDF := &MockPDFExtractor{Text: pdfText}
	mockLLM := &MockLLMClient{Response: jsonResp}

	extractor := NewExtractor(mockPDF, mockLLM)

	// Dummy reader
	reader := strings.NewReader("dummy pdf")

	data, err := extractor.Extract(context.Background(), reader, int64(reader.Len()))
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if data.FullName != "John Doe" {
		t.Errorf("expected John Doe, got %s", data.FullName)
	}
	if data.Company.Name != "Acme Corp" {
		t.Errorf("expected Acme Corp, got %s", data.Company.Name)
	}
}
