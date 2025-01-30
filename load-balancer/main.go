package main

import (
	"fmt"
	"load-balancer/configs"
	"load-balancer/loadbalancer"
	"load-balancer/server"
	"log"
)

func main() {
	cfg := configs.LoadConfig()

	var servers []server.Server
	for _, addr := range cfg.Servers {
		servers = append(servers, server.NewSimpleServer(addr, cfg.HealthCheck))
	}

	lb := loadbalancer.NewLoadBalancer(cfg.Port, servers)

	fmt.Printf("Load balancer running on port %s\n", cfg.Port)
	fmt.Printf("Proxying to %d servers:\n", len(servers))

	for _, s := range servers {
		fmt.Printf("- %s\n", s.Address())
	}

	if err := lb.Start(); err != nil {
		log.Fatalf("Failed to start load balancer: %v", err)
	}
}
