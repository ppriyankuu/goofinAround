package api

import (
	"concurrency-optimized-rate-limiter/internal/ratelimiter"
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	limiter *ratelimiter.RateLimiter
}

func NewHandler(limiter *ratelimiter.RateLimiter) *Handler {
	return &Handler{
		limiter: limiter,
	}
}

func (h *Handler) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Header.Get("X-Client-ID")
		if clientID == "" {
			http.Error(w, "Missing client ID", http.StatusBadRequest)
			return
		}

		allowed, err := h.limiter.Allow(r.Context(), clientID, 1)
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		if !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	total, rejected, duration := h.limiter.Metrics().SnapShot()
	avg := duration / time.Duration(total)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_requests":     total,
		"rejected_requests":  rejected,
		"average_processing": avg.String(),
	})
}
