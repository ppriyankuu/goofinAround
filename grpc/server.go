package main

import (
	"context"
	"fmt"
	"grpc/helloworld" // Import the generated package
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	helloworld.UnimplementedGreeterServer // Embed UnimplementedGreeterServer for future compatibility
}

func (s *server) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{
		Message: "Hello " + req.GetName(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{}) // Register the Greeter service
	reflection.Register(s)                         // Register reflection service for gRPC

	fmt.Println("Server listening on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
