package main

import (
	"concurrency-optimized-rate-limiter/api"
	"concurrency-optimized-rate-limiter/internal/ratelimiter"
	"log"
	"net/http"
)

func main() {
	// initializing rate limiter with 10 tokens/sec and 20 token capacity
	config := ratelimiter.NewConfig(10, 20)
	rl := ratelimiter.New(config)
	handler := api.NewHandler(rl)

	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", handler.MetricsHandler)
	mux.Handle("/", handler.RateLimitMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		}),
	))

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
