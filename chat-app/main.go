package main

import (
	"chat-app/config"
	"chat-app/database"
	pubsub "chat-app/pub-sub"
	"chat-app/routes"
	"chat-app/websocket"
	"log"
	"net/http"
)

func main() {
	config.LoadConfig()

	pubsub.InitRedis()
	database.InitDB()

	hub := websocket.NewHub(database.DB)
	go hub.Run()

	router := routes.InitRoutes(hub, database.DB)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
