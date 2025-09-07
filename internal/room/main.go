package room

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"open-meet/internal/util"

	"github.com/google/uuid"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type Room struct {
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRoomResponse struct {
	Room *Room `json:"room"`
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

	hostURL := os.Getenv("LIVEKIT_SERVER")
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	roomClient := lksdk.NewRoomServiceClient(
		hostURL,
		apiKey,
		apiSecret)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	room, err := roomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name:             generateRoomName(),
		RoomPreset:       "",
		EmptyTimeout:     0,
		DepartureTimeout: 0,
		MaxParticipants:  10,
		NodeId:           "",
		Metadata:         "",
		Egress:           nil,
		MinPlayoutDelay:  0,
		MaxPlayoutDelay:  0,
		SyncStreams:      false,
		ReplayEnabled:    false,
		Agents:           nil,
	})
	if err != nil {
		fmt.Printf("Error creating room: %v\n", err)
		http.Error(w, "Failed to create room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log room creation
	fmt.Printf("Room created - ID: %s Name: %s, Creator: %s, Time: %s\n",
		room.GetSid(),
		room.GetName(),
		userEmail,
		time.Now().Format(time.RFC3339))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := &CreateRoomResponse{
		Room: &Room{
			Name:      room.GetName(),
			CreatedBy: userEmail,
			CreatedAt: time.Now(),
		},
	}

	_ = json.NewEncoder(w).Encode(response)
}

// GetRoomHandler handles GET /rooms/{roomName} requests
func GetRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract room name from path
	roomName := r.URL.Path[len("/rooms/"):]
	if roomName == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}

	// Get LiveKit room client
	hostURL := os.Getenv("LIVEKIT_SERVER")
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")
	roomClient := lksdk.NewRoomServiceClient(hostURL, apiKey, apiSecret)

	// Get room from LiveKit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rooms, err := roomClient.ListRooms(ctx, &livekit.ListRoomsRequest{
		Names: []string{roomName},
	})
	if err != nil {
		http.Error(w, "Failed to list room, error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(rooms.GetRooms()) == 0 {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	room := rooms.GetRooms()[0]

	// Log room access
	fmt.Printf("Room accessed - Name: %s, Time: %s\n",
		room,
		time.Now().Format(time.RFC3339))

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"name":             room.Name,
		"num_participants": room.NumParticipants,
		"active_recording": room.ActiveRecording,
		"creation_time":    room.CreationTime,
		"sid":              room.Sid,
	})
}
