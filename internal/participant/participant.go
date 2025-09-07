package participant

import (
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
