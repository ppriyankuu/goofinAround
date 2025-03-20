package websockets

import (
	"chat-server/internals/models"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan models.Message)
	mutex     sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		conn.Close()
	}()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	go writePump()
	readPump(conn)
}

func readPump(conn *websocket.Conn) {
	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("ERRor reading message: %v", err)
			break
		}
		broadcast <- message
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	conn.Close()
}

func writePump() {
	for message := range broadcast {
		mutex.Lock()
		for client := range clients {
			if err := client.WriteJSON(message); err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
