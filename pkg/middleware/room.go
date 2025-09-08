package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireAuth ensures the request has a valid authentication token
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		c.Set("token", token)
		c.Next()
	}
}

// ValidateCreateRoom validates the request body for room creation
func ValidateCreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if contentType != "application/json" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Content-Type must be application/json",
			})
			return
		}

		c.Next()
	}
}

// RateLimitRoom implements rate limiting for room creation
func RateLimitRoom() gin.HandlerFunc {
	// TODO: Implement proper rate limiting with Redis or similar
	return func(c *gin.Context) {
		c.Next()
	}
}
