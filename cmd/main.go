package main

import (
	"fmt"
	"open-meet/pkg/api"
	"os"
)

func init() {
	// Validate required environment variables
	requiredEnvVars := []string{
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"OAUTH_REDIRECT_URL",
		"ALLOWED_ORIGINS",
		"LIVEKIT_API_KEY",
		"LIVEKIT_API_SECRET",
	}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			fmt.Printf("Error: %s environment variable is required\n", envVar)
			os.Exit(1)
		}
	}
}

func main() {
	cfg := &api.Config{
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		AllowedOrigins:     os.Getenv("ALLOWED_ORIGINS"),
		LiveKitServer:      os.Getenv("LIVEKIT_SERVER"),
		LiveKitAPIKey:      os.Getenv("LIVEKIT_API_KEY"),
		LiveKitAPISecret:   os.Getenv("LIVEKIT_API_SECRET"),
	}

	svc, err := api.NewEngine(cfg)
	if err != nil {
		panic(err)
	}

	err = svc.Run(":8080")
	if err != nil {
		panic(err)
	}
}
