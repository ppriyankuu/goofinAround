package main

import (
	"encoding/json"
	"log"
	"net/http"

	tollbooth "github.com/didip/tollbooth/v7"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

var defaultMessage = Message{
	Status: "Successful",
	Body:   "Hi! you've reached the API. How may I help you?",
}

func endpointHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(&defaultMessage); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	tlbthLimiter := tollbooth.NewLimiter(1, nil)
	tlbthLimiter.SetMessageContentType("application/json")
	tlbthLimiter.SetMessage(`{"status":"error","body":"Rate limit exceeded. Please try again later."}`)

	http.Handle("/ping", tollbooth.LimitFuncHandler(tlbthLimiter, endpointHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("There is something kind of an error the server is facing: %v", err)
	}
}
