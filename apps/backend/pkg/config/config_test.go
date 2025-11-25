package config_test

import (
	"os"
	"testing"

	"github.com/credfolio/apps/backend/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_WhenValidEnvVars_LoadsConfiguration(t *testing.T) {
	// Arrange
	t.Setenv("DATABASE_URL", "postgresql://test:test@localhost/test")
	t.Setenv("OPENAI_API_KEY", "sk-test-key")
	t.Setenv("PORT", "9000")
	t.Setenv("ENVIRONMENT", "testing")

	// Act
	cfg, err := config.Load()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "postgresql://test:test@localhost/test", cfg.DatabaseURL)
	assert.Equal(t, "sk-test-key", cfg.OpenAIAPIKey)
	assert.Equal(t, "9000", cfg.Port)
	assert.Equal(t, "testing", cfg.Environment)
}

func TestLoad_WhenMissingEnvVars_UsesDefaults(t *testing.T) {
	// Arrange - clear any existing env vars
	_ = os.Unsetenv("DATABASE_URL")
	_ = os.Unsetenv("OPENAI_API_KEY")
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("ENVIRONMENT")

	// Act
	cfg, err := config.Load()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, cfg.DatabaseURL, "localhost")
	assert.Empty(t, cfg.OpenAIAPIKey)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.Environment)
}

func TestConfig_IsDevelopment_WhenDevelopment_ReturnsTrue(t *testing.T) {
	// Arrange
	t.Setenv("ENVIRONMENT", "development")

	cfg, err := config.Load()
	require.NoError(t, err)

	// Act & Assert
	assert.True(t, cfg.IsDevelopment())
	assert.False(t, cfg.IsProduction())
}

func TestConfig_IsProduction_WhenProduction_ReturnsTrue(t *testing.T) {
	// Arrange
	t.Setenv("ENVIRONMENT", "production")

	cfg, err := config.Load()
	require.NoError(t, err)

	// Act & Assert
	assert.True(t, cfg.IsProduction())
	assert.False(t, cfg.IsDevelopment())
}
