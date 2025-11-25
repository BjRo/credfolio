package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseURL  string
	OpenAIAPIKey string
	Port         string
	Environment  string
	UploadPath   string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not present)
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:  getEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=credfolio port=5432 sslmode=disable"),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		Port:         getEnv("PORT", "8080"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		UploadPath:   getEnv("UPLOAD_PATH", "./uploads"),
	}

	return cfg, nil
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
