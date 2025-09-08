package api

import (
	"net/http"

	"open-meet/pkg/util"

	"github.com/gin-gonic/gin"
)

type LiveKitTokenRequest struct {
	RoomName string `json:"room_name" binding:"required"`
}

// LiveKitTokenHandler handles requests for generating LiveKit tokens
func (s *Service) LiveKitTokenHandler(c *gin.Context) {
	log := s.Log.WithName("livekit-token")

	req := new(LiveKitTokenRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		log.Error(err, "invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: room_name and identity are required",
		})
		return
	}
	if req.RoomName == "" {
		log.Info("room name or identity is empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: room_name and identity cannot be empty",
		})
		return
	}

	// Get user email from context (set by Authentication middleware)
	userEmail, err := util.GetUserEmailFromContext(c)
	if err != nil {
		log.Error(err, "failed to get user email from context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check if the room exists
	roomCtx := c.Request.Context()
	room, found, err := s.Store.Room().Get(roomCtx, req.RoomName)
	if err != nil {
		log.Error(err, "failed to get room info", "roomName", req.RoomName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if !found {
		log.Info("room not found", "roomName", req.RoomName)
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	// Generate token
	token, err := s.Store.Participant().GenerateToken(roomCtx, req.RoomName, userEmail)
	if err != nil {
		log.Error(err, "failed to generate token", "roomName", req.RoomName, "identity", userEmail)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
			"code":  "TOKEN_GENERATION_FAILED",
		})
		return
	}

	log.Info("token generated successfully", "roomName", req.RoomName, "identity", userEmail, "roomSid", room.Sid)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"room": gin.H{
			"name":             room.Name,
			"sid":              room.Sid,
			"num_participants": room.NumParticipants,
		},
	})
}
