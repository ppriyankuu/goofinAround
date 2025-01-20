package main

import (
	"context" // for context management
	"log"     // for logging
	"time"    // for setting timeouts

	pb "learn_grpc/greet" // Import the generated greet package

	"google.golang.org/grpc" // gRPC client library
)

func main() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dial the server using grpc.DialContext
	conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithInsecure()) // Without TLS for simplicity
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client for the GreetService
	client := pb.NewGreetServiceClient(conn)

	// Prepare a request
	req := &pb.GreetRequest{Name: "John"}

	// Call the Greet method
	respCtx, respCancel := context.WithTimeout(context.Background(), time.Second)
	defer respCancel()
	resp, err := client.Greet(respCtx, req)
	if err != nil {
		log.Fatalf("Error while calling Greet: %v", err)
	}

	// Print the response
	log.Printf("Response from server: %s", resp.GetMessage())
}
