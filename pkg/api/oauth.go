package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// CallbackHandler decodes the token to extract user information
func (s *Service) CallbackHandler(c *gin.Context) {
	log := s.Log.WithName("oauth-callback")

	var signInResponse GoogleSignInResponse
	if err := c.ShouldBindJSON(&signInResponse); err != nil {
		log.Error(err, "failed to decode request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if signInResponse.Credential == "" {
		log.Info("credential not provided in request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credential not provided in request body"})
		return
	}

	// Parse the JWT token without verification (already verified by middleware)
	token, _, err := new(jwt.Parser).ParseUnverified(signInResponse.Credential, &GoogleClaims{})
	if err != nil {
		log.Error(err, "failed to parse token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse token"})
		return
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		log.Error(nil, "failed to parse claims")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse claims"})
		return
	}

	// Log successful login
	log.Info("user information extracted",
		"email", claims.Email,
		"name", claims.Name)

	c.JSON(http.StatusOK, gin.H{
		"token":   signInResponse.Credential,
		"type":    "Bearer",
		"name":    claims.Name,
		"email":   claims.Email,
		"picture": claims.Picture,
	})
}
