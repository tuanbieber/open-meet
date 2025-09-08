package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetUserEmailFromContext extracts the authenticated user's email from gin context
func GetUserEmailFromContext(c *gin.Context) (string, error) {
	email, exists := c.Get("email")
	if !exists {
		return "", fmt.Errorf("email not found in context")
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", fmt.Errorf("email in context is not a string")
	}

	if emailStr == "" {
		return "", fmt.Errorf("email in context is empty")
	}

	return emailStr, nil
}
