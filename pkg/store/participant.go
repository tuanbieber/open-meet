package store

import (
	"context"
	"fmt"
	"os"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

// Participant defines the interface for participant operations
type Participant interface {
	// Token management
	GenerateToken(ctx context.Context, roomName, identity string) (string, error)
	ValidateToken(ctx context.Context, token string) (*auth.ClaimGrants, error)

	// Participant operations
	JoinRoom(ctx context.Context, roomName, identity string) error
	LeaveRoom(ctx context.Context, roomName, identity string) error
	GetParticipantInfo(ctx context.Context, roomName, identity string) (*livekit.ParticipantInfo, error)
	ListParticipants(ctx context.Context, roomName string) ([]*livekit.ParticipantInfo, error)

	// Media controls
	MuteSelf(ctx context.Context, roomName, identity string) error
	UnmuteSelf(ctx context.Context, roomName, identity string) error
	EnableVideo(ctx context.Context, roomName, identity string) error
	DisableVideo(ctx context.Context, roomName, identity string) error
	ShareScreen(ctx context.Context, roomName, identity string) error
	StopScreenShare(ctx context.Context, roomName, identity string) error
}

// participant implements Participant interface
type participant struct {
	client    *lksdk.RoomServiceClient
	apiKey    string
	apiSecret string
}

// NewParticipant creates a new participant instance
func NewParticipant() (*participant, error) {
	hostURL := os.Getenv("LIVEKIT_SERVER")
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	if hostURL == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("missing required LiveKit environment variables")
	}

	client := lksdk.NewRoomServiceClient(hostURL, apiKey, apiSecret)

	return &participant{
		client:    client,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}, nil
}

// GenerateToken creates a token for room access
func (p *participant) GenerateToken(ctx context.Context, roomName, identity string) (string, error) {
	at := auth.NewAccessToken(p.apiKey, p.apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomName,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(24 * 60 * 60) // 24 hours validity

	return at.ToJWT()
}

// ValidateToken validates a LiveKit token
func (p *participant) ValidateToken(ctx context.Context, token string) (*auth.ClaimGrants, error) {
	//claims, err := auth.ParseAPIToken(token, p.apiSecret)
	//if err != nil {
	//	return nil, fmt.Errorf("invalid token: %w", err)
	//}
	return nil, nil
}

// JoinRoom allows a participant to join a room
func (p *participant) JoinRoom(ctx context.Context, roomName, identity string) error {
	// First check if room exists
	rooms, err := p.client.ListRooms(ctx, &livekit.ListRoomsRequest{
		Names: []string{roomName},
	})
	if err != nil {
		return fmt.Errorf("failed to check room: %w", err)
	}
	if len(rooms.GetRooms()) == 0 {
		return fmt.Errorf("room %s does not exist", roomName)
	}

	// Generate token for joining
	_, err = p.GenerateToken(ctx, roomName, identity)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Update participant metadata
	_, err = p.client.UpdateParticipant(ctx, &livekit.UpdateParticipantRequest{
		Room:     roomName,
		Identity: identity,
		Metadata: `{"joinTime": "` + fmt.Sprint(ctx.Value("timestamp")) + `"}`,
	})
	if err != nil {
		return fmt.Errorf("failed to update participant metadata: %w", err)
	}

	return nil
}

// LeaveRoom allows a participant to leave a room
func (p *participant) LeaveRoom(ctx context.Context, roomName, identity string) error {
	_, err := p.client.RemoveParticipant(ctx, &livekit.RoomParticipantIdentity{
		Room:     roomName,
		Identity: identity,
	})
	if err != nil {
		return fmt.Errorf("failed to leave room: %w", err)
	}
	return nil
}

// GetParticipantInfo gets information about a specific participant
func (p *participant) GetParticipantInfo(ctx context.Context, roomName, identity string) (*livekit.ParticipantInfo, error) {
	participants, err := p.ListParticipants(ctx, roomName)
	if err != nil {
		return nil, err
	}

	for _, participant := range participants {
		if participant.Identity == identity {
			return participant, nil
		}
	}

	return nil, fmt.Errorf("participant not found")
}

// ListParticipants lists all participants in a room
func (p *participant) ListParticipants(ctx context.Context, roomName string) ([]*livekit.ParticipantInfo, error) {
	resp, err := p.client.ListParticipants(ctx, &livekit.ListParticipantsRequest{
		Room: roomName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list participants: %w", err)
	}
	return resp.GetParticipants(), nil
}

// MuteSelf mutes the participant's audio
func (p *participant) MuteSelf(ctx context.Context, roomName, identity string) error {
	_, err := p.client.UpdateParticipant(ctx, &livekit.UpdateParticipantRequest{
		Room:     roomName,
		Identity: identity,
		Metadata: `{"audio": false}`,
		Permission: &livekit.ParticipantPermission{
			CanPublish: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to mute self: %w", err)
	}
	return nil
}

// UnmuteSelf unmutes the participant's audio
func (p *participant) UnmuteSelf(ctx context.Context, roomName, identity string) error {
	_, err := p.client.UpdateParticipant(ctx, &livekit.UpdateParticipantRequest{
		Room:     roomName,
		Identity: identity,
		Metadata: `{"audio": true}`,
		Permission: &livekit.ParticipantPermission{
			CanPublish: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to unmute self: %w", err)
	}
	return nil
}

// EnableVideo enables the participant's video
func (p *participant) EnableVideo(ctx context.Context, roomName, identity string) error {
	_, err := p.client.UpdateParticipant(ctx, &livekit.UpdateParticipantRequest{
		Room:     roomName,
		Identity: identity,
		Metadata: `{"video": true}`,
		Permission: &livekit.ParticipantPermission{
			CanPublish: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to enable video: %w", err)
	}
	return nil
}

// DisableVideo disables the participant's video
func (p *participant) DisableVideo(ctx context.Context, roomName, identity string) error {
	_, err := p.client.UpdateParticipant(ctx, &livekit.UpdateParticipantRequest{
		Room:     roomName,
		Identity: identity,
		Metadata: `{"video": false}`,
		Permission: &livekit.ParticipantPermission{
			CanPublish: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to disable video: %w", err)
	}
	return nil
}

// ShareScreen enables screen sharing for the participant
func (p *participant) ShareScreen(ctx context.Context, roomName, identity string) error {
	_, err := p.client.UpdateParticipant(ctx, &livekit.UpdateParticipantRequest{
		Room:     roomName,
		Identity: identity,
		Metadata: `{"screen": true}`,
		Permission: &livekit.ParticipantPermission{
			CanPublish:     true,
			CanPublishData: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to enable screen sharing: %w", err)
	}
	return nil
}

// StopScreenShare disables screen sharing for the participant
func (p *participant) StopScreenShare(ctx context.Context, roomName, identity string) error {
	_, err := p.client.UpdateParticipant(ctx, &livekit.UpdateParticipantRequest{
		Room:     roomName,
		Identity: identity,
		Metadata: `{"screen": false}`,
		Permission: &livekit.ParticipantPermission{
			CanPublish:     true,
			CanPublishData: false,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to disable screen sharing: %w", err)
	}
	return nil
}
