package services

import (
	"chat-app/internal/models"
	pubsub "chat-app/pkg/pub-sub"
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

type RoomService struct {
	rooms          map[string]*models.Room
	mu             sync.Mutex
	Pubsub         *pubsub.RedisPubSub
	messageService *MessageService
}

func NewRoomService(redisClient *redis.Client, messageService *MessageService) *RoomService {
	return &RoomService{
		rooms:          make(map[string]*models.Room),
		Pubsub:         pubsub.NewRedisPubSub(redisClient),
		messageService: messageService,
	}
}

// Helper method to get a room safely
func (rs *RoomService) GetRoom(roomID string) (*models.Room, bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	room, exists := rs.rooms[roomID]
	return room, exists
}

// Helper method to add a user to a room safely
func (rs *RoomService) addUserToRoom(roomID, userID string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	room, exists := rs.rooms[roomID]
	if !exists {
		room = &models.Room{
			ID:       roomID,
			Users:    make(map[string]bool),
			MaxUsers: 5,
		}
		rs.rooms[roomID] = room
	}

	if len(room.Users) >= room.MaxUsers {
		return errors.New("room is full")
	}

	room.Users[userID] = true
	return nil
}

// allows a user to join a room
func (rs *RoomService) JoinRoom(roomID, userID string) error {
	if roomID == "" || userID == "" {
		return errors.New("roomID and userID must not be empty")
	}

	if err := rs.addUserToRoom(roomID, userID); err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	return nil
}

// LeaveRoom allows a user to leave a room
func (rs *RoomService) LeaveRoom(roomID, userID string) {
	if roomID == "" || userID == "" {
		return
	}

	rs.mu.Lock()
	defer rs.mu.Unlock()

	if room, exists := rs.rooms[roomID]; exists {
		delete(room.Users, userID)
		if len(room.Users) == 0 {
			delete(rs.rooms, roomID)
		}
	}
}

// broadcasts a message to all users in a room
func (rs *RoomService) BroadcastMessage(roomID, userID, content string) error {
	if roomID == "" || userID == "" || content == "" {
		return errors.New("roomID, userID, and content must not be empty")
	}

	// Creating the message object
	message := &models.Message{
		RoomID:  roomID,
		UserID:  userID,
		Content: content,
	}

	// Saving the message using the message service
	if err := rs.messageService.SaveMessage(message); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	// Publishing the message to the room's pub/sub channel
	if err := rs.Pubsub.Publish(roomID, message); err != nil {
		return fmt.Errorf("failed to broadcast message: %w", err)
	}

	return nil
}
