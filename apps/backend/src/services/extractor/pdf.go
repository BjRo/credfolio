package extractor

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ledongthuc/pdf"
)

type PDFExtractor struct{}

func NewPDFExtractor() *PDFExtractor {
	return &PDFExtractor{}
}

// ExtractText reads text from a PDF reader
func (e *PDFExtractor) ExtractText(r io.ReaderAt, size int64) (string, error) {
	f, err := pdf.NewReader(r, size)
	if err != nil {
		return "", fmt.Errorf("failed to open pdf: %w", err)
	}

	var buf bytes.Buffer
	// Iterate over all pages
	numPages := f.NumPage()
	for i := 1; i <= numPages; i++ {
		p := f.Page(i)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			// Log warning? For now just continue
			continue
		}
		buf.WriteString(text)
		buf.WriteString("\n")
	}

	return buf.String(), nil
}
