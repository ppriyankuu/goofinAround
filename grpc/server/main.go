package main

import (
	"context" // for context management
	"log"     // for logging
	"net"     // for networking

	pb "learn_grpc/greet" // Import the generated greet package

	"google.golang.org/grpc" // gRPC server library
)

// Server struct
type Server struct {
	pb.UnimplementedGreetServiceServer // Embed for forward compatibility
}

// Greet implements the GreetService
func (s *Server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("Received: %s", req.GetName())
	return &pb.GreetResponse{Message: "Hello, " + req.GetName()}, nil
}

func main() {
	// Create a listener on TCP port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Register the GreetService with the server
	pb.RegisterGreetServiceServer(grpcServer, &Server{})

	log.Println("Server is listening on port 50051...")
	// Start the server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
