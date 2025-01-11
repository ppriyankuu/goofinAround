package routes

import (
	"chat-app/models"
	"chat-app/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(hub *websocket.Hub, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// WebSocket endpoint with room query parameter
	router.GET("/ws", func(c *gin.Context) {
		_ = c.DefaultQuery("room", "general01") // default to "general01" if no room is provided
		// Pass the room in the URL query string directly to the WebSocket handler
		websocket.ServeWebSocket(hub, c.Writer, c.Request)
	})

	// Messages endpoint for retrieving messages in a specific room
	router.GET("/messages/:room", func(c *gin.Context) {
		room := c.Param("room")
		messages, err := models.GetMessages(db, room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}
		c.JSON(http.StatusOK, messages)
	})
	return router
}
