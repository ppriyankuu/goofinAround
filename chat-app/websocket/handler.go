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

	// create a new client
	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte),
		Room: "general",
	}

	// register the client with the hub
	hub.Register <- client

	// start go routines for reading and writing messages
	go client.ReadPump()
	go client.WritePump()

}
