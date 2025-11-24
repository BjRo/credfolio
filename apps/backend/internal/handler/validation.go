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

	// Check length
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

	// Check for only whitespace
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

// SanitizeString sanitizes a string by:
// - Trimming whitespace
// - Escaping HTML entities
// - Removing control characters (except newlines, tabs, carriage returns)
func SanitizeString(value string) string {
	// Trim whitespace
	sanitized := strings.TrimSpace(value)

	// Remove control characters except newlines, tabs, and carriage returns
	var builder strings.Builder
	for _, r := range sanitized {
		if r == '\n' || r == '\t' || r == '\r' {
			builder.WriteRune(r)
		} else if r >= 32 || r == '\n' || r == '\t' || r == '\r' {
			builder.WriteRune(r)
		}
		// Skip other control characters
	}

	// HTML escape to prevent XSS
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

	// Get base name to prevent path traversal
	baseName := filepath.Base(filename)

	// Remove any remaining path separators
	baseName = strings.ReplaceAll(baseName, "/", "")
	baseName = strings.ReplaceAll(baseName, "\\", "")

	// Validate length
	if len(baseName) > 255 {
		return "", &ValidationError{
			ErrorCode: ErrorCodeInvalidRequest,
			Message:   "filename is too long (max 255 characters)",
			Field:     "filename",
		}
	}

	// Check for invalid characters
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

// ValidateProfileSummary validates a profile summary input
func ValidateProfileSummary(value *string) error {
	// Summary is optional, but if provided, should be valid
	if value == nil {
		return nil
	}
	return ValidateString(value, "summary", 0, 2000, false)
}

// SanitizeProfileSummary sanitizes a profile summary
func SanitizeProfileSummary(value string) string {
	return SanitizeString(value)
}
