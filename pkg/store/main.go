package store

import (
	"sync"
)

var (
	roomStore        Room
	roomStoreOnce    sync.Once
	hostStore        Host
	hostStoreOnce    sync.Once
	participantStore Participant
	participantOnce  sync.Once
)

// GetHostStore returns the singleton Host instance
func GetHostStore() (Host, error) {
	var initErr error
	hostStoreOnce.Do(func() {
		var store *host
		store, initErr = NewHost()
		if initErr == nil {
			hostStore = store
		}
	})
	if initErr != nil {
		return nil, initErr
	}
	return hostStore, nil
}

// GetRoomStore returns the singleton Room instance
func GetRoomStore() (Room, error) {
	var initErr error
	roomStoreOnce.Do(func() {
		var store *LiveKitRoom
		store, initErr = NewLiveKitRoom()
		if initErr == nil {
			roomStore = store
		}
	})
	if initErr != nil {
		return nil, initErr
	}
	return roomStore, nil
}

// GetParticipantStore returns the singleton Participant instance
func GetParticipantStore() (Participant, error) {
	var initErr error
	participantOnce.Do(func() {
		var store *participant
		store, initErr = NewParticipant()
		if initErr == nil {
			participantStore = store
		}
	})
	if initErr != nil {
		return nil, initErr
	}
	return participantStore, nil
}

// Store represents the main data store interface
type Store interface {
	Room() Room
	Host() Host
	Participant() Participant
}

// memoryStore implements Store interface
type memoryStore struct {
	room        Room
	host        Host
	participant Participant
}

func NewStore() (Store, error) {
	roomSt, err := GetRoomStore()
	if err != nil {
		return nil, err
	}

	hostSt, err := GetHostStore()
	if err != nil {
		return nil, err
	}

	participantSt, err := GetParticipantStore()
	if err != nil {
		return nil, err
	}

	return &memoryStore{
		room:        roomSt,
		host:        hostSt,
		participant: participantSt,
	}, nil
}

func (s *memoryStore) Room() Room {
	return s.room
}

func (s *memoryStore) Host() Host {
	return s.host
}

func (s *memoryStore) Participant() Participant {
	return s.participant
}
