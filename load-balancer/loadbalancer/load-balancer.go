package loadbalancer

import (
	"fmt"
	"load-balancer/server"
	"net/http"
	"sync/atomic"
)

type LoadBalancer struct {
	port    string
	counter uint64
	servers []server.Server
}

func NewLoadBalancer(port string, servers []server.Server) *LoadBalancer {
	return &LoadBalancer{
		port:    port,
		servers: servers,
	}
}

func (lb *LoadBalancer) GetNextAvailableServer() (server.Server, error) {
	numServers := uint64(len(lb.servers))
	attempts := numServers

	for i := uint64(0); i < attempts; i++ {
		idx := atomic.AddUint64(&lb.counter, 1) % numServers
		server := lb.servers[idx]
		if server.IsAlive() {
			return server, nil
		}
	}

	return nil, fmt.Errorf("no server available")
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server, err := lb.GetNextAvailableServer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	fmt.Printf("Forwarding request to: %s\n", server.Address())
	server.Serve(w, r)
}

func (lb *LoadBalancer) Start() error {
	return http.ListenAndServe(":"+lb.port, lb)
}
