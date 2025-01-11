package websocket

import (
	"chat-app/models"
	"log"

	"gorm.io/gorm"
)

// Hub manages all WebSocket clients and broadcasts messages to them
type Hub struct {
	Rooms      map[string]map[*Client]bool // Rooms with connected clients
	Clients    map[*Client]bool            // Connected clients
	Broadcast  chan Message                // Broadcast channel for messages
	Register   chan *Client                // Register new clients
	Unregister chan *Client                // Unregister clients
	DB         *gorm.DB                    // Database connection
}

// Message structure to include room info
type Message struct {
	Room    string
	Content []byte
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		DB:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// Ensure Rooms map is initialized
			if h.Rooms == nil {
				h.Rooms = make(map[string]map[*Client]bool)
			}

			// Register client in the specified room
			if h.Rooms[client.Room] == nil {
				h.Rooms[client.Room] = make(map[*Client]bool)
			}
			h.Rooms[client.Room][client] = true

		case client := <-h.Unregister:
			// Unregister client from the specified room
			if room, ok := h.Rooms[client.Room]; ok {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.Send)

					// clean up empty rooms
					if len(room) == 0 {
						delete(h.Rooms, client.Room)
					}
				}
			}

		case message := <-h.Broadcast:
			// Ensure the Rooms map is initialized
			if h.Rooms == nil {
				h.Rooms = make(map[string]map[*Client]bool)
			}

			// Save message to the database
			roomID := getRoomID(message.Room) // Fetch the room ID based on the room name
			err := models.SaveMessage(h.DB, string(message.Content), "user", message.Room, roomID)
			if err != nil {
				log.Printf("Failed to save message: %v", err)
			}

			// Broadcast to all clients in the specified room
			if clients, ok := h.Rooms[message.Room]; ok {
				for client := range clients {
					select {
					case client.Send <- message.Content:
					default:
						delete(clients, client)
						close(client.Send)
					}
				}
			}
		}
	}
}

// Dummy function to fetch room ID based on room name (adjust as needed)
func getRoomID(roomName string) uint {
	// Return the room ID based on the room name
	if roomName == "general01" {
		return 1
	} else if roomName == "general02" {
		return 2
	} else if roomName == "general03" {
		return 3
	}
	return 0 // Return 0 for unknown rooms
}
