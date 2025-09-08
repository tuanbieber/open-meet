package store

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

// Room defines the interface for room operations
type Room interface {
	Create(ctx context.Context, name, creatorEmail string) (*livekit.Room, error)
	Get(ctx context.Context, name string) (*livekit.Room, bool, error)
	List(ctx context.Context) ([]*livekit.Room, error)
	Delete(ctx context.Context, name string) error
	SetHost(roomName, hostEmail string)
	GetRoomHost(roomName string) (string, bool)
	IsHost(roomName, email string) bool
}

// LiveKitRoom implements Room interface
type LiveKitRoom struct {
	client *lksdk.RoomServiceClient
	mu     sync.RWMutex
	hosts  map[string]string // map[roomName]hostEmail
}

func NewLiveKitRoom() (*LiveKitRoom, error) {
	hostURL := os.Getenv("LIVEKIT_SERVER")
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	if hostURL == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("missing required LiveKit environment variables")
	}

	client := lksdk.NewRoomServiceClient(hostURL, apiKey, apiSecret)

	return &LiveKitRoom{
		client: client,
		hosts:  make(map[string]string),
	}, nil
}

func (r *LiveKitRoom) Create(ctx context.Context, name, creatorEmail string) (*livekit.Room, error) {
	room, err := r.client.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name:             name,
		EmptyTimeout:     30 * 60, // 30 minutes
		DepartureTimeout: 5 * 60,  // 5 minutes
		MaxParticipants:  100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	// Set creator as host
	r.SetHost(name, creatorEmail)
	return room, nil
}

func (r *LiveKitRoom) Get(ctx context.Context, name string) (*livekit.Room, bool, error) {
	resp, err := r.client.ListRooms(ctx, &livekit.ListRoomsRequest{
		Names: []string{name},
	})
	if err != nil {
		return nil, false, fmt.Errorf("failed to get room %s: %w", name, err)
	}

	rooms := resp.GetRooms()
	if len(rooms) == 0 {
		return nil, false, nil
	}

	return rooms[0], true, nil
}

// List returns all rooms from LiveKit server
func (r *LiveKitRoom) List(ctx context.Context) ([]*livekit.Room, error) {
	resp, err := r.client.ListRooms(ctx, &livekit.ListRoomsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list rooms: %w", err)
	}
	return resp.GetRooms(), nil
}

// Delete removes a room from LiveKit server
func (r *LiveKitRoom) Delete(ctx context.Context, name string) error {
	_, err := r.client.DeleteRoom(ctx, &livekit.DeleteRoomRequest{
		Room: name,
	})
	if err != nil {
		return fmt.Errorf("failed to delete room %s: %w", name, err)
	}

	// Remove host mapping
	r.mu.Lock()
	delete(r.hosts, name)
	r.mu.Unlock()

	return nil
}

// SetHost sets the host for a room
func (r *LiveKitRoom) SetHost(roomName, hostEmail string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hosts[roomName] = hostEmail
}

// GetRoomHost returns the host email for a room
func (r *LiveKitRoom) GetRoomHost(roomName string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	host, exists := r.hosts[roomName]
	return host, exists
}

// IsHost checks if the given email is the host of the room
func (r *LiveKitRoom) IsHost(roomName, email string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	host, exists := r.hosts[roomName]
	return exists && host == email
}
