package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port        string
	DatabaseURL string
	OpenAIKey   string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		OpenAIKey:   getEnv("OPENAI_API_KEY", ""),
	}
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
