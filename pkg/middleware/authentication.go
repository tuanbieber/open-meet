package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

// Authentication validates Google Sign-In JWT tokens and checks request content type
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Content-Type for POST requests
		if c.Request.Method == http.MethodPost {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "Content-Type must be application/json",
				})
				return
			}
		}

		// Check Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization header provided",
			})
			return
		}

		// Strip "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// Validate Google Sign-In JWT token
		payload, err := validateGoogleToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("invalid token: %v", err),
			})
			return
		}

		// Verify token was issued for our application
		if payload.Audience != os.Getenv("GOOGLE_CLIENT_ID") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token has wrong audience",
			})
			return
		}

		// Verify email is verified by Google
		if !payload.Claims["email_verified"].(bool) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "email not verified by Google",
			})
			return
		}

		// Store validated token data in context
		c.Set("token", token)
		c.Set("email", payload.Claims["email"])
		c.Set("name", payload.Claims["name"])
		c.Set("picture", payload.Claims["picture"])

		c.Next()
	}
}

// validateGoogleToken validates a Google Sign-In JWT token
func validateGoogleToken(ctx context.Context, tokenString string) (*idtoken.Payload, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID environment variable not set")
	}

	// Verify the token using Google's public keys
	payload, err := idtoken.Validate(ctx, tokenString, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	return payload, nil
}
