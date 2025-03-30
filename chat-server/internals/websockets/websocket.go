package websockets

import (
	"chat-server/internals/models"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 'clients' maintains a thread-safe set of active WebSocket connections.
var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan models.Message) // Tracks connected clients.
	mutex     sync.Mutex                  // Protects concurrent access to clients.

	// upgrader configures the WebSocket upgrade process.
	upgrader = websocket.Upgrader{
		// Allow any origin for simplicity.
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// HandleWebSocket manages WebSocket connection upgrades and lifecycle.
func HandleWebSocket(c *gin.Context) {
	// Upgrade HTTP request to WebSocket protocol.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Ensure proper closure of connection.
	defer func() {
		mutex.Lock()
		delete(clients, conn) // Remove client on disconnect.
		mutex.Unlock()
		conn.Close()
	}()

	// Register new client.
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	// Start a separate goroutine to handle message broadcasting.
	go writePump()

	// Read messages from the client.
	readPump(conn)
}

// readPump continuously reads messages from a WebSocket connection.
func readPump(conn *websocket.Conn) {
	for {
		var message models.Message

		// Read JSON message from the WebSocket.
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Send the received message to the broadcast channel.
		broadcast <- message
	}

	// Cleanup when the loop breaks (client disconnected)
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	conn.Close()
}

// writePump continuously listen for messages on the broadcast channel
// and sends them to all clients.
func writePump() {
	for message := range broadcast {
		// Locking to ensure thread-safe iteration.
		mutex.Lock()
		for client := range clients {
			// Write message to the client
			if err := client.WriteJSON(message); err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(clients, client) // Remove client if writing fails.
			}
		}
		mutex.Unlock()
	}
}
