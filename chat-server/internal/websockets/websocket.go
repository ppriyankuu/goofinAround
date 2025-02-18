package websockets

import (
	"chat-server/internal/models"
	"chat-server/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development purposes
	},
}

type WebSocketServer struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan models.Message
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	service    *services.ChatService
}

func NewWebSocketServer(service *services.ChatService) *WebSocketServer {
	return &WebSocketServer{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan models.Message),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		service:    service,
	}
}

func (wss *WebSocketServer) Run() {
	for {
		select {
		case client := <-wss.register:
			wss.clients[client] = true
		case client := <-wss.unregister:
			if _, ok := wss.clients[client]; ok {
				client.Close()
				delete(wss.clients, client)
			}
		case message := <-wss.broadcast:
			for client := range wss.clients {
				err := client.WriteJSON(message)
				if err != nil {
					client.Close()
					delete(wss.clients, client)
				}
			}
		}
	}
}

func (wss *WebSocketServer) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer conn.Close()

	wss.register <- conn
	defer func() { wss.unregister <- conn }()

	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		savedMessage, err := wss.service.SendMessage(message.UserID, message.GroupID, message.Content)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}

		wss.broadcast <- *savedMessage
	}
}
