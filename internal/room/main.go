package room

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"open-meet/internal/util"

	"github.com/google/uuid"
)

type Room struct {
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRoomResponse struct {
	Room *Room `json:"room"`
}

// generateRoomID creates a shorter room ID by hashing a UUID
func generateRoomID() string {
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

func NewRoom(createdBy string) *Room {
	return &Room{
		Name:      generateRoomID(),
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
}

// CreateRoomHandler handles POST /rooms requests
func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user email from token
	userEmail, err := util.GetUserEmailFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Create new room with hashed ID
	newRoom := NewRoom(userEmail) // name parameter is ignored now
	response := CreateRoomResponse{Room: newRoom}

	// Log room creation
	fmt.Printf("Room created - Name: %s, Creator: %s, Time: %s\n",
		newRoom.Name,
		newRoom.CreatedBy,
		newRoom.CreatedAt.Format(time.RFC3339))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
