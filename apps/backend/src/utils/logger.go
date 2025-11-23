package utils

import (
	"log/slog"
	"os"
)

// SetupLogger configures the default slog logger.
// It returns the logger instance.
func SetupLogger() *slog.Logger {
	// Default to text handler for dev, json for prod could be an option
	// For now, using JSON handler as it's structured
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
