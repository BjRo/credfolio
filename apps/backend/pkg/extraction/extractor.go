package extraction

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extractor handles text extraction from various file formats
type Extractor struct{}

// NewExtractor creates a new text extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractFromFile extracts text content from a file based on its extension
func (e *Extractor) ExtractFromFile(filePath string) (string, error) {
	if filePath == "" {
		return "", errors.New("file path is required")
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".txt":
		return e.extractFromTextFile(filePath)
	case ".md", ".markdown":
		return e.extractFromMarkdownFile(filePath)
	default:
		return "", errors.New("unsupported file type: " + ext)
	}
}

// ExtractFromReader extracts text content from a reader
func (e *Extractor) ExtractFromReader(reader io.Reader, fileType string) (string, error) {
	if reader == nil {
		return "", errors.New("reader is required")
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// extractFromTextFile extracts text from a plain text file
func (e *Extractor) extractFromTextFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// extractFromMarkdownFile extracts text from a markdown file
// For now, we just return the raw markdown content
// The LLM can understand markdown format
func (e *Extractor) extractFromMarkdownFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// SupportedExtensions returns a list of supported file extensions
func (e *Extractor) SupportedExtensions() []string {
	return []string{".txt", ".md", ".markdown"}
}

// IsSupportedExtension checks if a file extension is supported
func (e *Extractor) IsSupportedExtension(ext string) bool {
	ext = strings.ToLower(ext)
	for _, supported := range e.SupportedExtensions() {
		if ext == supported {
			return true
		}
	}
	return false
}
