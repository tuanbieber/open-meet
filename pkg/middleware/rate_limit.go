package middleware

import (
	"github.com/gin-gonic/gin"
)

// RateLimit implements rate limiting for room creation
func RateLimit() gin.HandlerFunc {
	// TODO: Implement proper rate limiting with Redis or similar
	return func(c *gin.Context) {
		c.Next()
	}
}
