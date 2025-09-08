package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

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

// CallbackHandler processes Google Sign-In callback with token verification
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
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

	// Verify the token with Google
	payload, err := verifyGoogleToken(r.Context(), signInResponse.Credential)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Print user email to server logs
	fmt.Printf("User logged in - Email: %s\n", payload.Claims["email"])

	// Return user information as JSON response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"token":   signInResponse.Credential,
		"type":    "Bearer",
		"name":    payload.Claims["name"],
		"email":   payload.Claims["email"],
		"picture": payload.Claims["picture"],
	})
}

// verifyGoogleToken verifies that the token is actually from Google
func verifyGoogleToken(ctx context.Context, tokenString string) (*idtoken.Payload, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID environment variable not set")
	}

	payload, err := idtoken.Validate(ctx, tokenString, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}

	// Verify that the token was issued for our application
	if payload.Audience != clientID {
		return nil, fmt.Errorf("token has wrong audience: %s", payload.Audience)
	}

	// Additional security checks
	if !payload.Claims["email_verified"].(bool) {
		return nil, fmt.Errorf("email not verified by Google")
	}

	return payload, nil
}
