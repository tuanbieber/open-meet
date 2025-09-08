package api

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"open-meet/pkg/util"
)

// CreateRoomHandler handles POST /rooms requests
func (s *Service) CreateRoomHandler(c *gin.Context) {
	log := s.Log.WithName("create-room")

	// Get user email from context (set by Authentication middleware)
	userEmail, err := util.GetUserEmailFromContext(c)
	if err != nil {
		log.Error(err, "failed to get user email from context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	roomName := generateRoomName()
	lkRoom, err := s.Store.Room().Create(c.Request.Context(), roomName, userEmail)
	if err != nil {
		log.Error(err, "failed to create room",
			"roomName", roomName,
			"creator", userEmail)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log room creation
	log.Info("room created",
		"roomID", lkRoom.GetSid(),
		"roomName", lkRoom.GetName(),
		"creator", userEmail)

	response := &CreateRoomResponse{
		Room: &Room{
			Name:      lkRoom.GetName(),
			CreatedBy: userEmail,
			CreatedAt: time.Now(),
		},
	}

	c.JSON(http.StatusCreated, response)
}

// GetRoomHandler handles GET /rooms/:roomName requests
func (s *Service) GetRoomHandler(c *gin.Context) {
	log := s.Log.WithName("get-room")

	roomName := c.Param("room_name")
	if roomName == "" {
		log.Info("room name not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "room name is required"})
		return
	}

	lkRoom, found, err := s.Store.Room().Get(c.Request.Context(), roomName)
	if err != nil {
		log.Error(err, "failed to get room", "roomName", roomName)
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	if !found {
		log.Info("room not found", "roomName", roomName)
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	log.Info("room accessed",
		"roomName", lkRoom.GetName(),
		"numParticipants", lkRoom.NumParticipants)

	c.JSON(http.StatusOK, gin.H{
		"name":             lkRoom.Name,
		"num_participants": lkRoom.NumParticipants,
		"active_recording": lkRoom.ActiveRecording,
		"creation_time":    lkRoom.CreationTime,
		"sid":              lkRoom.Sid,
	})
}

func generateRoomName() string {
	// Generate a UUID first
	newString := uuid.NewString()

	// Create SHA-256 hash
	hasher := sha256.New()
	hasher.Write([]byte(newString))
	hash := hasher.Sum(nil)

	// Convert to URL-safe base64 and take first 10 characters
	encoded := base64.URLEncoding.EncodeToString(hash)
	return encoded[:10]
}

type Room struct {
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRoomResponse struct {
	Room *Room `json:"room"`
}
