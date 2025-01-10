package main

import (
	"context"
	"fmt"
	"grpc/helloworld" // Import the generated package
	"log"

	"google.golang.org/grpc"
)

func main() {
	// Connect to the server
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new Greeter client
	c := helloworld.NewGreeterClient(conn)

	// Call SayHello RPC
	resp, err := c.SayHello(context.Background(), &helloworld.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Println("Greeting:", resp.GetMessage())
}
