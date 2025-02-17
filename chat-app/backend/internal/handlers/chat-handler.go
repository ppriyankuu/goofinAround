package handlers

import (
	"chat-app/internal/models"
	"chat-app/internal/services"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type ChatHandler struct {
	roomService    *services.RoomService
	messageService *services.MessageService
}

func NewChatHandler(roomService *services.RoomService, messageService *services.MessageService) *ChatHandler {
	return &ChatHandler{
		roomService:    roomService,
		messageService: messageService,
	}
}

func (ch *ChatHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ws/:roomID/:userID", ch.handleWebSocket)
	r.GET("/rooms/:roomID/messages", ch.getMessages)
}

func (ch *ChatHandler) handleWebSocket(c *gin.Context) {
	roomID := c.Param("roomID")
	userID := c.Param("userID")

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Join the room
	if err := ch.roomService.JoinRoom(roomID, userID); err != nil {
		log.Printf("Failed to join room: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer ch.roomService.LeaveRoom(roomID, userID)

	// Subscribe to the room's pubsub channel
	sub, err := ch.roomService.Pubsub.Subscribe(roomID)
	if err != nil {
		log.Printf("Failed to subscribe to room: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer sub.Close()

	// Channel to signal when the connection should close
	stopChan := make(chan struct{})

	// Goroutine to handle incoming messages from the client
	go func() {
		defer close(stopChan)
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Failed to read message: %v", err)
				return
			}

			var msg models.Message
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return
			}

			if err := ch.roomService.BroadcastMessage(roomID, userID, msg.Content); err != nil {
				log.Printf("Failed to broadcast message: %v", err)
				return
			}
		}
	}()

	// Goroutine to handle outgoing messages to the client
	go func() {
		for {
			select {
			case <-stopChan:
				return
			default:
				msg, err := sub.Receive(context.Background())
				if err != nil {
					log.Printf("Failed to receive message: %v", err)
					return
				}

				if err := conn.WriteJSON(msg); err != nil {
					log.Printf("Failed to write message: %v", err)
					return
				}
			}
		}
	}()

	// Wait for the stop signal
	<-stopChan
}

func (ch *ChatHandler) getMessages(c *gin.Context) {
	roomID := c.Param("roomID")

	// Parse the limit query parameter
	limitStr := c.Query("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 10 // Default limit if not provided or invalid
	}

	// Fetch messages for the specified room
	messages, err := ch.messageService.GetMessagesByRoom(roomID)
	if err != nil {
		log.Printf("Failed to fetch messages for room '%s': %v", roomID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	// Return the messages as JSON
	c.JSON(http.StatusOK, messages)
}
