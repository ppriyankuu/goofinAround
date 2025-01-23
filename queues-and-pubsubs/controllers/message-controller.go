package controllers

import (
	"log"
	"net/http"
	"queues-and-pubsubs/utils"

	"github.com/gin-gonic/gin"
)

func PublishMessage(c *gin.Context) {
	topic := c.Param("topic")
	var message struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err := utils.RDB.Publish(utils.CTX, topic, message.Message).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Message published"})
}

func SubscribeToTopic(c *gin.Context) {
	topic := c.Param("topic")
	subscriber := utils.RDB.Subscribe(utils.CTX, topic)
	ch := subscriber.Channel()

	go func() {
		for message := range ch {
			log.Printf("Received message from topic %s: %s", topic, message.Payload)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"status": "Subscribed to topic", "topic": topic})
}
