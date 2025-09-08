package main

import (
	"fmt"
	"open-meet/pkg/api"
	"os"

	"github.com/gin-gonic/gin"
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

	gin.SetMode(gin.DebugMode)

	err = svc.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func backup() {
	//mux := http.NewServeMux()
	//
	//mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "application/json")
	//	_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	//})
	//
	//mux.HandleFunc("/livekit-tokens", participant.LiveKitTokenHandler)
	//mux.HandleFunc("/callback", oauth.CallbackHandler)
	//mux.HandleFunc("/rooms", room.CreateRoomHandler)
	//mux.HandleFunc("/rooms/", room.GetRoomHandler)
	//
	//// Setup CORS
	//c := cors.New(cors.Options{
	//	AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGINS")},
	//	AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
	//	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	//	AllowCredentials: true,
	//	Debug:            false,
	//})
	//
	//// Wrap the mux with CORS middleware
	//handler := c.Handler(mux)
	//
	//fmt.Println("Server starting on :8080")
	//_ = http.ListenAndServe(":8080", handler)
}
