package service

import (
	"fmt"

	"github.com/credfolio/apps/backend/pkg/logger"
)

// BaseService provides common functionality for all services
type BaseService struct {
	logger *logger.Logger
}

// NewBaseService creates a new base service with logger
func NewBaseService(logger *logger.Logger) *BaseService {
	return &BaseService{
		logger: logger,
	}
}

// Logger returns the logger instance
func (b *BaseService) Logger() *logger.Logger {
	return b.logger
}

// WrapError wraps an error with a consistent format
// Usage: return nil, b.WrapError(err, "failed to create profile")
func (b *BaseService) WrapError(err error, operation string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", operation, err)
}

// WrapErrorf wraps an error with a formatted message
// Usage: return nil, b.WrapErrorf(err, "failed to create profile for user %s", userID)
func (b *BaseService) WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	operation := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", operation, err)
}

// ValidateNotEmpty validates that a string is not empty after trimming whitespace
// Returns an error if the string is empty
func (b *BaseService) ValidateNotEmpty(value, fieldName string) error {
	if value == "" || len([]rune(value)) == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	// Check if string is only whitespace
	trimmed := ""
	for _, r := range value {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			trimmed += string(r)
		}
	}
	if trimmed == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// LogOperationStart logs the start of an operation
func (b *BaseService) LogOperationStart(operation string, args ...interface{}) {
	b.logger.Info(operation, args...)
}

// LogOperationComplete logs the completion of an operation
func (b *BaseService) LogOperationComplete(operation string, args ...interface{}) {
	b.logger.Info(operation, args...)
}

// LogError logs an error with context
func (b *BaseService) LogError(message string, args ...interface{}) {
	b.logger.Error(message, args...)
}
