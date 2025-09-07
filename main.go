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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig *oauth2.Config
)

func init() {
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

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

	switch r.Method {
	case http.MethodPost:
		var signInResponse GoogleSignInResponse
		if err := json.NewDecoder(r.Body).Decode(&signInResponse); err != nil {
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
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token":   signInResponse.Credential,
			"type":    "Bearer",
			"name":    claims.Name,
			"email":   claims.Email,
			"picture": claims.Picture,
		})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
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
