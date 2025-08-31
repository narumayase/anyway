package config

import (
	"github.com/joho/godotenv"
	anysherlog "github.com/narumayase/anysher/log"
	"github.com/rs/zerolog/log"
	"os"
)

// Config contains the application configuration
type Config struct {
	Port string
}

// Load loads configuration from environment variables or an .env file
func Load() Config {
	anysherlog.SetLogLevel()
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Debug().Msgf("No .env file found or error loading .env file: %v", err)
	}
	return Config{
		Port: getEnv("PORT", "8080"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
