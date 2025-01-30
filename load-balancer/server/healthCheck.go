package server

import (
	"net/http"
	"time"
)

func (s *SimpleServer) runHealthCheck() {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for {
		<-time.After(s.healthCheck)
		response, err := client.Head(s.addr)
		if err != nil || response.StatusCode >= 500 {
			s.markDown()
		} else {
			s.markUp()
		}
	}
}

func (s *SimpleServer) markUp() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alive = true
}
func (s *SimpleServer) markDown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alive = false
}
