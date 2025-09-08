package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port           string
	AllowedOrigins string

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	OAuthRedirectURL   string

	// LiveKit
	LiveKitServer    string
	LiveKitAPIKey    string
	LiveKitAPISecret string
}

// LoadConfig loads environment variables from .env file and returns Config
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// List of required environment variables
	required := []string{
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"OAUTH_REDIRECT_URL",
		"ALLOWED_ORIGINS",
		"LIVEKIT_API_KEY",
		"LIVEKIT_API_SECRET",
		"LIVEKIT_SERVER",
	}

	// Check for missing environment variables
	for _, env := range required {
		if os.Getenv(env) == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", env)
		}
	}

	return &Config{
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		OAuthRedirectURL:   os.Getenv("OAUTH_REDIRECT_URL"),
		LiveKitServer:      os.Getenv("LIVEKIT_SERVER"),
		LiveKitAPIKey:      os.Getenv("LIVEKIT_API_KEY"),
		LiveKitAPISecret:   os.Getenv("LIVEKIT_API_SECRET"),
		AllowedOrigins:     os.Getenv("ALLOWED_ORIGINS"),
		Port:               os.Getenv("PORT"),
	}, nil
}
