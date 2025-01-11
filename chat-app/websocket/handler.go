package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWebSocket handles WebSocket connections and upgrades HTTP to WebSocket
func ServeWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Get the room parameter from the URL	query string
	room := r.URL.Query().Get("room") // Fetch room from query string
	if room == "" {
		room = "general01" // Default room if not provided
	}

	// create a new client
	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte),
		Room: room,
	}

	// register the client with the hub
	hub.Register <- client

	// start go routines for reading and writing messages
	go client.ReadPump()
	go client.WritePump()

}
