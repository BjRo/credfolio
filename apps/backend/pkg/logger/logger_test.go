package logger_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/credfolio/apps/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestLogger_Info_WhenLoggingInfoMessage_WritesToOutput(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	l := logger.NewWithOutput(logger.LevelInfo, &buf)

	// Act
	l.Info("test message")

	// Assert
	output := buf.String()
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "test message")
}

func TestLogger_Debug_WhenLevelIsInfo_DoesNotWrite(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	l := logger.NewWithOutput(logger.LevelInfo, &buf)

	// Act
	l.Debug("debug message")

	// Assert
	assert.Empty(t, buf.String())
}

func TestLogger_Debug_WhenLevelIsDebug_WritesToOutput(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	l := logger.NewWithOutput(logger.LevelDebug, &buf)

	// Act
	l.Debug("debug message")

	// Assert
	output := buf.String()
	assert.Contains(t, output, "[DEBUG]")
	assert.Contains(t, output, "debug message")
}

func TestLogger_Warn_WritesToOutput(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	l := logger.NewWithOutput(logger.LevelInfo, &buf)

	// Act
	l.Warn("warning message")

	// Assert
	output := buf.String()
	assert.Contains(t, output, "[WARN]")
	assert.Contains(t, output, "warning message")
}

func TestLogger_Error_WritesToOutput(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	l := logger.NewWithOutput(logger.LevelInfo, &buf)

	// Act
	l.Error("error message")

	// Assert
	output := buf.String()
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "error message")
}

func TestLogger_Info_WithFormatArgs_FormatsMessage(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	l := logger.NewWithOutput(logger.LevelInfo, &buf)

	// Act
	l.Info("user %s logged in with id %d", "john", 123)

	// Assert
	output := buf.String()
	assert.Contains(t, output, "user john logged in with id 123")
}

func TestLogger_SetLevel_ChangesLogLevel(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	l := logger.NewWithOutput(logger.LevelError, &buf)

	// Initially should not log info
	l.Info("should not appear")
	assert.Empty(t, buf.String())

	// Act
	l.SetLevel(logger.LevelInfo)
	l.Info("should appear")

	// Assert
	assert.True(t, strings.Contains(buf.String(), "should appear"))
}

func TestNew_CreatesLoggerWithDefaultOutput(t *testing.T) {
	// Act
	l := logger.New(logger.LevelInfo)

	// Assert
	assert.NotNil(t, l)
}
