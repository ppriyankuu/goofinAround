package websocket

import (
	"chat-app/models"
	"log"

	"gorm.io/gorm"
)

// Hub manages all WebSocket clients and broadcasts messages to them
type Hub struct {
	Clients    map[*Client]bool // Connected clients
	Broadcast  chan []byte      // Broadcast channel for messages
	Register   chan *Client     // Register new clients
	Unregister chan *Client     // Unregister clients
	DB         *gorm.DB         // Database connection
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		DB:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			// Save message to the database with roomID
			roomID := getRoomID("general") // Function to fetch room ID by name (adjust this based on your actual setup)
			err := models.SaveMessage(h.DB, string(message), "user", "general", roomID)
			if err != nil {
				log.Printf("Failed to save message: %v", err)
			}

			// Send message to all clients in the room
			for client := range h.Clients {
				if client.Room == "general" { // Broadcasting to the "general" room for now
					select {
					case client.Send <- message:
					default:
						delete(h.Clients, client)
						close(client.Send)
					}
				}
			}
		}
	}
}

// Dummy function to fetch room ID (replace with actual logic)
func getRoomID(_ string) uint {
	// Lookup or calculate the roomID based on roomName
	return 1 // Placeholder value, replace with actual logic
}
