package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(http.ResponseWriter, *http.Request)
}

type SimpleServer struct {
	addr        string
	proxy       *httputil.ReverseProxy
	mu          sync.RWMutex
	alive       bool
	healthCheck time.Duration
}

func NewSimpleServer(addr string, healthCheckInterval time.Duration) *SimpleServer {
	serverURL, err := url.Parse(addr)
	if err != nil {
		log.Fatalf("Invalid server address: %v", err)
	}

	s := &SimpleServer{
		addr:        addr,
		proxy:       httputil.NewSingleHostReverseProxy(serverURL),
		healthCheck: healthCheckInterval,
		alive:       true,
	}

	go s.runHealthCheck()
	return s
}

func (s *SimpleServer) Address() string { return s.addr }

func (s *SimpleServer) IsAlive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.alive
}

func (s *SimpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}
