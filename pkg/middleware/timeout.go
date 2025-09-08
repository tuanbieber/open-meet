package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout middleware wraps the request with a timeout
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Update the request with the new context
		c.Request = c.Request.WithContext(ctx)

		// Create a channel to signal request completion
		done := make(chan struct{})
		go func() {
			// Execute the next handlers
			c.Next()
			done <- struct{}{}
		}()

		// Wait for either timeout or completion
		select {
		case <-done:
			// Request completed before timeout
			return
		case <-ctx.Done():
			// Timeout exceeded
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "Request timeout exceeded",
				"code":  "REQUEST_TIMEOUT",
			})
			return
		}
	}
}
