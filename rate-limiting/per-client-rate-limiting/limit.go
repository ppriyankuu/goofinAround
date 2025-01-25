package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func perClientRateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// cleanup goroutine to remove inactive clients
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		if ip == "" {
			http.Error(w, "Unable to determine client IP", http.StatusInternalServerError)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			message := Message{
				Status: "Request failed",
				Body:   "You are sending too many requests. Slow down!",
			}

			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&message)
			return

		}
		mu.Unlock()
		// pass the control to the next function (controller)
		next(w, r)
	})
}

func getIP(r *http.Request) string {
	// Check for headers set by reverse proxies
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}
