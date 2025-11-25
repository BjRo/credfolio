package extraction_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/credfolio/apps/backend/pkg/extraction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractor_ExtractFromFile_WhenGivenTxtFile_ReturnsTextContent(t *testing.T) {
	// Arrange
	extractor := extraction.NewExtractor()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "This is a test reference letter content."
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	// Act
	result, err := extractor.ExtractFromFile(filePath)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, content, result)
}

func TestExtractor_ExtractFromFile_WhenGivenMarkdownFile_ReturnsTextContent(t *testing.T) {
	// Arrange
	extractor := extraction.NewExtractor()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")
	content := "# Reference Letter\n\nThis is **bold** content."
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	// Act
	result, err := extractor.ExtractFromFile(filePath)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, content, result)
}

func TestExtractor_ExtractFromFile_WhenEmptyPath_ReturnsError(t *testing.T) {
	// Arrange
	extractor := extraction.NewExtractor()

	// Act
	result, err := extractor.ExtractFromFile("")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "file path is required")
}

func TestExtractor_ExtractFromFile_WhenUnsupportedExtension_ReturnsError(t *testing.T) {
	// Arrange
	extractor := extraction.NewExtractor()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.pdf")
	err := os.WriteFile(filePath, []byte("dummy"), 0644)
	require.NoError(t, err)

	// Act
	result, err := extractor.ExtractFromFile(filePath)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "unsupported file type")
}

func TestExtractor_ExtractFromFile_WhenFileNotFound_ReturnsError(t *testing.T) {
	// Arrange
	extractor := extraction.NewExtractor()

	// Act
	result, err := extractor.ExtractFromFile("/nonexistent/path/file.txt")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestExtractor_ExtractFromReader_WhenGivenReader_ReturnsContent(t *testing.T) {
	// Arrange
	extractor := extraction.NewExtractor()
	content := "This is reader content"
	reader := strings.NewReader(content)

	// Act
	result, err := extractor.ExtractFromReader(reader, "txt")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, content, result)
}

func TestExtractor_ExtractFromReader_WhenNilReader_ReturnsError(t *testing.T) {
	// Arrange
	extractor := extraction.NewExtractor()

	// Act
	result, err := extractor.ExtractFromReader(nil, "txt")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "reader is required")
}

func TestExtractor_IsSupportedExtension_WhenTxt_ReturnsTrue(t *testing.T) {
	// Arrange
	extractor := extraction.NewExtractor()

	// Act & Assert
	assert.True(t, extractor.IsSupportedExtension(".txt"))
	assert.True(t, extractor.IsSupportedExtension(".TXT"))
}

func TestExtractor_IsSupportedExtension_WhenMarkdown_ReturnsTrue(t *testing.T) {
	// Arrange
	extractor := extraction.NewExtractor()

	// Act & Assert
	assert.True(t, extractor.IsSupportedExtension(".md"))
	assert.True(t, extractor.IsSupportedExtension(".markdown"))
}

func TestExtractor_IsSupportedExtension_WhenUnsupported_ReturnsFalse(t *testing.T) {
	// Arrange
	extractor := extraction.NewExtractor()

	// Act & Assert
	assert.False(t, extractor.IsSupportedExtension(".pdf"))
	assert.False(t, extractor.IsSupportedExtension(".docx"))
}
