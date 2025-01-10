package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	url := "http://localhost:8080/secure" // Replace with your endpoint

	// Define the authorization header
	authHeader := "Bearer crankyyy"

	// Number of requests to send
	requestCount := 10

	for i := 1; i <= requestCount; i++ {
		// Create a new HTTP request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Request %d failed: %s\n", i, err)
			continue
		}

		// Add the Authorization header
		req.Header.Set("Authorization", authHeader)

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Request %d failed: %s\n", i, err)
			continue
		}

		// Read the response body
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Print the status code and response
		fmt.Printf("Request %d: Status %d, Response: %s\n", i, resp.StatusCode, strings.TrimSpace(string(body)))

		// Wait for a second before sending the next request
		time.Sleep(1 * time.Second)
	}
}
