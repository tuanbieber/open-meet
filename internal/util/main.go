package util

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// GoogleClaims represents the claims in a Google JWT token
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	jwt.RegisteredClaims
}

// GetUserEmailFromToken extracts JWT token from Authorization header and gets user email
func GetUserEmailFromToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}

	// Remove "Bearer " prefix if present
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	// Parse the JWT token
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &GoogleClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse claims")
	}

	return claims.Email, nil
}
