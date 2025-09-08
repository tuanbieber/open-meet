package store

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

// Host defines the interface for host operations
type Host interface {
	// Room management
	EndMeeting(ctx context.Context, roomName string, hostEmail string) error
	LockRoom(ctx context.Context, roomName string, hostEmail string) error
	UnlockRoom(ctx context.Context, roomName string, hostEmail string) error

	// Participant management
	KickParticipant(ctx context.Context, roomName string, hostEmail string, participantIdentity string) error
	MuteParticipant(ctx context.Context, roomName string, hostEmail string, participantIdentity string) error
	UnmuteParticipant(ctx context.Context, roomName string, hostEmail string, participantIdentity string) error

	// Host management
	AssignHost(ctx context.Context, roomName string, currentHostEmail string, newHostEmail string) error
	IsHost(roomName, email string) bool
	GetRoomHost(roomName string) (string, bool)
}

// host implements Host interface
type host struct {
	client *lksdk.RoomServiceClient
	mu     sync.RWMutex
	hosts  map[string]string // map[roomName]hostEmail
}

// NewHost creates a new host instance
func NewHost() (*host, error) {
	hostURL := os.Getenv("LIVEKIT_SERVER")
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	if hostURL == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("missing required LiveKit environment variables")
	}

	client := lksdk.NewRoomServiceClient(hostURL, apiKey, apiSecret)

	return &host{
		client: client,
		hosts:  make(map[string]string),
	}, nil
}

// EndMeeting terminates the meeting for all participants
func (h *host) EndMeeting(ctx context.Context, roomName string, hostEmail string) error {
	if !h.IsHost(roomName, hostEmail) {
		return fmt.Errorf("unauthorized: only host can end meeting")
	}

	_, err := h.client.DeleteRoom(ctx, &livekit.DeleteRoomRequest{
		Room: roomName,
	})
	if err != nil {
		return fmt.Errorf("failed to end meeting: %w", err)
	}

	// Remove host mapping
	h.mu.Lock()
	delete(h.hosts, roomName)
	h.mu.Unlock()

	return nil
}

// LockRoom prevents new participants from joining
func (h *host) LockRoom(ctx context.Context, roomName string, hostEmail string) error {
	if !h.IsHost(roomName, hostEmail) {
		return fmt.Errorf("unauthorized: only host can lock room")
	}

	_, err := h.client.UpdateRoomMetadata(ctx, &livekit.UpdateRoomMetadataRequest{
		Room:     roomName,
		Metadata: `{"locked": true}`,
	})
	if err != nil {
		return fmt.Errorf("failed to lock room: %w", err)
	}

	return nil
}

// UnlockRoom allows new participants to join
func (h *host) UnlockRoom(ctx context.Context, roomName string, hostEmail string) error {
	if !h.IsHost(roomName, hostEmail) {
		return fmt.Errorf("unauthorized: only host can unlock room")
	}

	_, err := h.client.UpdateRoomMetadata(ctx, &livekit.UpdateRoomMetadataRequest{
		Room:     roomName,
		Metadata: `{"locked": false}`,
	})
	if err != nil {
		return fmt.Errorf("failed to unlock room: %w", err)
	}

	return nil
}

// KickParticipant removes a participant from the room
func (h *host) KickParticipant(ctx context.Context, roomName string, hostEmail string, participantIdentity string) error {
	if !h.IsHost(roomName, hostEmail) {
		return fmt.Errorf("unauthorized: only host can kick participants")
	}

	// Prevent host from kicking themselves
	if participantIdentity == hostEmail {
		return fmt.Errorf("host cannot kick themselves")
	}

	_, err := h.client.RemoveParticipant(ctx, &livekit.RoomParticipantIdentity{
		Room:     roomName,
		Identity: participantIdentity,
	})
	if err != nil {
		return fmt.Errorf("failed to kick participant: %w", err)
	}

	return nil
}

// MuteParticipant disables a participant's audio
func (h *host) MuteParticipant(ctx context.Context, roomName string, hostEmail string, participantIdentity string) error {
	if !h.IsHost(roomName, hostEmail) {
		return fmt.Errorf("unauthorized: only host can mute participants")
	}

	// Get the participant's tracks
	_, err := h.client.UpdateParticipant(ctx, &livekit.UpdateParticipantRequest{
		Room:     roomName,
		Identity: participantIdentity,
		Metadata: `{"audio": false}`,
		Permission: &livekit.ParticipantPermission{
			CanPublish: false,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to mute participant: %w", err)
	}

	return nil
}

// UnmuteParticipant enables a participant's audio
func (h *host) UnmuteParticipant(ctx context.Context, roomName string, hostEmail string, participantIdentity string) error {
	if !h.IsHost(roomName, hostEmail) {
		return fmt.Errorf("unauthorized: only host can unmute participants")
	}

	_, err := h.client.UpdateParticipant(ctx, &livekit.UpdateParticipantRequest{
		Room:     roomName,
		Identity: participantIdentity,
		Metadata: `{"audio": true}`,
		Permission: &livekit.ParticipantPermission{
			CanPublish: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to unmute participant: %w", err)
	}

	return nil
}

// AssignHost transfers host privileges to another participant
func (h *host) AssignHost(ctx context.Context, roomName string, currentHostEmail string, newHostEmail string) error {
	if !h.IsHost(roomName, currentHostEmail) {
		return fmt.Errorf("unauthorized: only current host can assign new host")
	}

	h.mu.Lock()
	h.hosts[roomName] = newHostEmail
	h.mu.Unlock()

	return nil
}

// IsHost checks if the given email is the host of the room
func (h *host) IsHost(roomName, email string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	host, exists := h.hosts[roomName]
	return exists && host == email
}

// GetRoomHost returns the host email for a room
func (h *host) GetRoomHost(roomName string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	host, exists := h.hosts[roomName]
	return host, exists
}
