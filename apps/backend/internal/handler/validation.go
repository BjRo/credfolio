package handler

import (
	"fmt"
	"html"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

// ValidationError represents a validation error with a specific error code
type ValidationError struct {
	ErrorCode int
	Message   string
	Field     string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field %s: %s", e.Field, e.Message)
}

// ValidateString validates a string input with length constraints
func ValidateString(value *string, fieldName string, minLen, maxLen int, required bool) error {
	if value == nil {
		if required {
			return &ValidationError{
				ErrorCode: ErrorCodeMissingField,
				Message:   fmt.Sprintf("%s is required", fieldName),
				Field:     fieldName,
			}
		}
		return nil
	}

	str := *value

	length := utf8.RuneCountInString(str)
	if length < minLen {
		return &ValidationError{
			ErrorCode: ErrorCodeInvalidRequest,
			Message:   fmt.Sprintf("%s must be at least %d characters", fieldName, minLen),
			Field:     fieldName,
		}
	}
	if length > maxLen {
		return &ValidationError{
			ErrorCode: ErrorCodeInvalidRequest,
			Message:   fmt.Sprintf("%s must be at most %d characters", fieldName, maxLen),
			Field:     fieldName,
		}
	}

	trimmed := strings.TrimSpace(str)
	if required && len(trimmed) == 0 {
		return &ValidationError{
			ErrorCode: ErrorCodeInvalidRequest,
			Message:   fmt.Sprintf("%s cannot be empty or only whitespace", fieldName),
			Field:     fieldName,
		}
	}

	return nil
}

func SanitizeString(value string) string {
	sanitized := strings.TrimSpace(value)

	var builder strings.Builder
	for _, r := range sanitized {
		if r == '\n' || r == '\t' || r == '\r' {
			builder.WriteRune(r)
		} else if r >= 32 || r == '\n' || r == '\t' || r == '\r' {
			builder.WriteRune(r)
		}
	}

	return html.EscapeString(builder.String())
}

// SanitizeFilename sanitizes a filename to prevent path traversal attacks
func SanitizeFilename(filename string) (string, error) {
	if filename == "" {
		return "", &ValidationError{
			ErrorCode: ErrorCodeInvalidRequest,
			Message:   "filename cannot be empty",
			Field:     "filename",
		}
	}

	baseName := filepath.Base(filename)
	baseName = strings.ReplaceAll(baseName, "/", "")
	baseName = strings.ReplaceAll(baseName, "\\", "")

	if len(baseName) > 255 {
		return "", &ValidationError{
			ErrorCode: ErrorCodeInvalidRequest,
			Message:   "filename is too long (max 255 characters)",
			Field:     "filename",
		}
	}

	invalidChars := []string{"..", "<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(baseName, char) {
			return "", &ValidationError{
				ErrorCode: ErrorCodeInvalidRequest,
				Message:   fmt.Sprintf("filename contains invalid character: %s", char),
				Field:     "filename",
			}
		}
	}

	return baseName, nil
}

// ValidateUUID validates a UUID string
func ValidateUUID(value string, fieldName string) (uuid.UUID, error) {
	if value == "" {
		return uuid.Nil, &ValidationError{
			ErrorCode: ErrorCodeMissingField,
			Message:   fmt.Sprintf("%s is required", fieldName),
			Field:     fieldName,
		}
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, &ValidationError{
			ErrorCode: ErrorCodeInvalidRequest,
			Message:   fmt.Sprintf("%s is not a valid UUID", fieldName),
			Field:     fieldName,
		}
	}

	return parsed, nil
}

// ValidateJobDescription validates a job description input
func ValidateJobDescription(value *string) error {
	return ValidateString(value, "jobDescription", 10, 10000, true)
}

// SanitizeJobDescription sanitizes a job description
func SanitizeJobDescription(value string) string {
	return SanitizeString(value)
}

func ValidateProfileSummary(value *string) error {
	if value == nil {
		return nil
	}
	return ValidateString(value, "summary", 0, 2000, false)
}

// SanitizeProfileSummary sanitizes a profile summary
func SanitizeProfileSummary(value string) string {
	return SanitizeString(value)
}

func ValidateFileType(mimeType string, filename string) error {
	allowedMimeTypes := []string{
		"text/plain",
		"text/markdown",
		"text/md",
		"text/txt",
		"text/x-md",
		"text/x-txt",
		"text/x-markdown",
		"text/x-md",
		"text/x-txt",
	}

	mimeTypeLower := strings.ToLower(strings.TrimSpace(mimeType))
	isValidMimeType := false
	for _, allowed := range allowedMimeTypes {
		if mimeTypeLower == allowed {
			isValidMimeType = true
			break
		}
	}

	ext := strings.ToLower(filepath.Ext(filename))
	isValidExtension := ext == ".md" || ext == ".txt" || ext == ".markdown"

	if !isValidMimeType && !isValidExtension {
		return &ValidationError{
			ErrorCode: ErrorCodeInvalidFileType,
			Message:   "file must be a txt or md file",
			Field:     "file",
		}
	}

	return nil
}

func ValidateFileSize(fileSize int64, maxSizeBytes int64) error {
	if fileSize <= 0 {
		return &ValidationError{
			ErrorCode: ErrorCodeInvalidRequest,
			Message:   "file is empty",
			Field:     "file",
		}
	}

	if fileSize > maxSizeBytes {
		maxSizeMB := float64(maxSizeBytes) / (1024 * 1024)
		fileSizeMB := float64(fileSize) / (1024 * 1024)
		return &ValidationError{
			ErrorCode: ErrorCodeFileTooLarge,
			Message:   fmt.Sprintf("file size (%.2f MB) exceeds maximum allowed size (%.2f MB)", fileSizeMB, maxSizeMB),
			Field:     "file",
		}
	}

	return nil
}
