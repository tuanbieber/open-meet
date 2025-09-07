package participant

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/livekit/protocol/auth"
)

const (
	// Token validity duration
	tokenDuration = 1 * time.Hour
)

func GenerateLiveKitToken(roomName, identity string) (string, error) {
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		Room:     roomName,
		RoomJoin: true,
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

// LiveKitTokenHandler handles POST /get-livekit-token requests
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
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
