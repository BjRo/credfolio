package pdf

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ledongthuc/pdf"
)

// ExtractorInterface defines the interface for PDF extraction
type ExtractorInterface interface {
	ExtractText(reader io.Reader) (string, error)
}

// Extractor handles PDF text extraction
type Extractor struct{}

// NewExtractor creates a new PDF extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Ensure Extractor implements ExtractorInterface
var _ ExtractorInterface = (*Extractor)(nil)

// ExtractText extracts text from a PDF reader
// Note: ledongthuc/pdf requires io.ReaderAt, so we need to read the entire file first
func (e *Extractor) ExtractText(reader io.Reader) (string, error) {
	// Read all data into memory since pdf.NewReader requires io.ReaderAt
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read PDF data: %w", err)
	}

	// Create a reader from bytes
	readerAt := bytes.NewReader(data)

	pdfReader, err := pdf.NewReader(readerAt, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var text string
	numPages := pdfReader.NumPage()
	for i := 1; i <= numPages; i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		pageText, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %w", i, err)
		}
		text += pageText + "\n"
	}

	return text, nil
}
