package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/cors"
)

//var (
//	oauthConfig *oauth2.Config
//)

func init() {
	//oauthConfig = &oauth2.Config{
	//	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	//	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	//	RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
	//	Scopes: []string{
	//		"https://www.googleapis.com/auth/userinfo.profile",
	//		"https://www.googleapis.com/auth/userinfo.email",
	//	},
	//	Endpoint: google.Endpoint,
	//}

	// Validate required environment variables
	requiredEnvVars := []string{"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET", "OAUTH_REDIRECT_URL", "ALLOWED_ORIGINS"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			fmt.Printf("Error: %s environment variable is required\n", envVar)
			os.Exit(1)
		}
	}
}

type GoogleSignInResponse struct {
	Credential string `json:"credential"`
	ClientID   string `json:"clientId"`
	SelectBy   string `json:"select_by"`
}

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	jwt.RegisteredClaims
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	// Log request information
	//fmt.Printf("Method: %s\n", r.Method)
	//fmt.Printf("Headers: %+v\n", r.Header)

	// Read and log the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	// Restore the body for further processing
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	signInResponse := new(GoogleSignInResponse)
	err = json.NewDecoder(r.Body).Decode(signInResponse)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if signInResponse.Credential == "" {
		http.Error(w, "Credential not provided in request body", http.StatusBadRequest)
		return
	}

	// Parse the JWT token without verification since it's already verified by Google
	token, _, err := new(jwt.Parser).ParseUnverified(signInResponse.Credential, &GoogleClaims{})
	if err != nil {
		http.Error(w, "Failed to parse token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
		return
	}

	// Print user email to server logs
	fmt.Printf("User logged in - Email: %s\n", claims.Email)

	// Return user information as JSON response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"token":   signInResponse.Credential,
		"type":    "Bearer",
		"name":    claims.Name,
		"email":   claims.Email,
		"picture": claims.Picture,
	})

}

type LiveKitTokenRequest struct {
	RoomName string `json:"room_name"`
	Identity string `json:"identity"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	mux.HandleFunc("/get-livekit-token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LiveKitTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		if req.RoomName == "" || req.Identity == "" {
			http.Error(w, "room and identity are required in request body", http.StatusBadRequest)
			return
		}

		// Generate LiveKit token
		token, err := GenerateLiveKitToken(req.RoomName, req.Identity)
		if err != nil {
			http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	})

	mux.HandleFunc("/callback", handleCallback)

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
