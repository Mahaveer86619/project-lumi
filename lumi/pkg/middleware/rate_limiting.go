package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DefaultLimit        = 10
	DefaultWindow       = 1 * time.Minute
	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderXRealIP       = "X-Real-Ip"
)

type ClientData struct {
	count int
	reset time.Time
	mu    sync.Mutex
}

type RateLimiter struct {
	clients map[string]*ClientData
	limit   int
	window  time.Duration
	mu      sync.Mutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if window <= 0 {
		window = DefaultWindow
	}

	return &RateLimiter{
		clients: make(map[string]*ClientData),
		limit:   limit,
		window:  window,
	}
}

func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get(HeaderXForwardedFor); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	if realIP := r.Header.Get(HeaderXRealIP); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		rl.mu.Lock()
		client, exists := rl.clients[ip]
		if !exists {
			client = &ClientData{
				count: 0,
				reset: time.Now().Add(rl.window),
			}
			rl.clients[ip] = client
		}
		rl.mu.Unlock()

		client.mu.Lock()
		defer client.mu.Unlock()

		if time.Now().After(client.reset) {
			client.count = 0
			client.reset = time.Now().Add(rl.window)
			log.Printf("Rate limit window reset for IP: %s. New window ends at: %v", ip, client.reset.Format(time.Kitchen))
		}

		if client.count >= rl.limit {
			log.Printf("Rate limit exceeded for IP: %s. Limit: %d/%s", ip, rl.limit, rl.window)

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(client.reset.Unix(), 10))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)

			waitTime := time.Until(client.reset).Round(time.Second).String()

			response := map[string]interface{}{
				"status_code": http.StatusTooManyRequests,
				"message":     fmt.Sprintf("Too Many Requests. Try again in %s", waitTime),
			}

			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Printf("Error writing JSON response: %v", err)
			}
			return
		}

		client.count++

		// Set informational headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(rl.limit-client.count))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(client.reset.Unix(), 10))

		next.ServeHTTP(w, r)
	}
}
