package extractor

import (
	"bytes"
	"testing"
)

func TestExtractText_InvalidPDF(t *testing.T) {
	extractor := NewPDFExtractor()
	data := []byte("not a pdf file")
	reader := bytes.NewReader(data)

	_, err := extractor.ExtractText(reader, int64(len(data)))
	if err == nil {
		t.Error("expected error for invalid pdf, got nil")
	}
}
