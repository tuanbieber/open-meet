package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"open-meet/internal/oauth"
	"open-meet/internal/participant"
	"open-meet/internal/room"

	"github.com/rs/cors"
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
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	mux.HandleFunc("/livekit-tokens", participant.LiveKitTokenHandler)
	mux.HandleFunc("/callback", oauth.CallbackHandler)
	mux.HandleFunc("/rooms", room.CreateRoomHandler)
	mux.HandleFunc("/rooms/", room.GetRoomHandler)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGINS")},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false,
	})

	// Wrap the mux with CORS middleware
	handler := c.Handler(mux)

	fmt.Println("Server starting on :8080")
	_ = http.ListenAndServe(":8080", handler)
}
