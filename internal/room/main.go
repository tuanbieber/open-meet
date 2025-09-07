package room

import (
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

func NewRoom(name, createdBy string) *Room {
	return &Room{
		Name:      name,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
}

// CreateHandler handles POST /rooms requests
func CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	// Use authenticated user's email as creator
	newRoomName := uuid.NewString()
	newRoom := NewRoom(newRoomName, userEmail)
	response := CreateRoomResponse{Room: newRoom}

	// Log room creation
	fmt.Printf("Room created - Name: %s, Creator: %s, Time: %s\n",
		newRoom.Name,
		newRoom.CreatedBy,
		newRoom.CreatedAt.Format(time.RFC3339))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
