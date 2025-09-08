package participant

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

const (
	// Token validity duration
	tokenDuration = 1 * time.Hour
)

func GenerateLiveKitToken(roomName, identity string) (string, error) {
	// Check if room exists using LiveKit's room service
	hostURL := os.Getenv("LIVEKIT_SERVER")
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	// Create room service client
	roomClient := lksdk.NewRoomServiceClient(hostURL, apiKey, apiSecret)

	// Check if room exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rooms, err := roomClient.ListRooms(ctx, &livekit.ListRoomsRequest{
		Names: []string{roomName},
	})
	if err != nil {
		return "", fmt.Errorf("failed to get room, error: %v", err)
	}
	if len(rooms.Rooms) == 0 {
		return "", fmt.Errorf("room %s does not exist", roomName)
	}

	// Generate token for existing room
	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomCreate:           false,
		RoomList:             false,
		RoomRecord:           false,
		RoomAdmin:            false,
		RoomJoin:             true,
		Room:                 roomName,
		CanPublish:           nil,
		CanSubscribe:         nil,
		CanPublishData:       nil,
		CanPublishSources:    nil,
		CanUpdateOwnMetadata: nil,
		IngressAdmin:         false,
		Hidden:               false,
		Recorder:             false,
		Agent:                false,
		CanSubscribeMetrics:  nil,
		DestinationRoom:      "",
	}
	at.SetVideoGrant(grant).
		SetIdentity(identity).
		SetValidFor(tokenDuration)

	// TODO: How long does the token last?

	return at.ToJWT()
}

type LiveKitTokenRequest struct {
	RoomName string `json:"room_name"`
	Identity string `json:"identity"`
}

func LiveKitTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LiveKitTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.RoomName == "" || req.Identity == "" {
		http.Error(w, "room and identity are required in request body", http.StatusBadRequest)
		return
	}

	// Generate LiveKit token
	token, err := GenerateLiveKitToken(req.RoomName, req.Identity)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
