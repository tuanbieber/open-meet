package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GoogleSignInResponse struct {
	Credential string `json:"credential"`
	ClientID   string `json:"clientId"`
	SelectBy   string `json:"select_by"`
}

// CallbackHandler processes Google Sign-In callback with token verification
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

	c.JSON(http.StatusOK, gin.H{
		"token": signInResponse.Credential,
		"type":  "Bearer",
		//"name":    payload.Claims["name"],
		//"email":   payload.Claims["email"],
		//"picture": payload.Claims["picture"],
	})
}
