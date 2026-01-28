package middleware

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Limit per IP
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients         = make(map[string]*client)
	mu              sync.Mutex
	rateLimit       = rate.Every(time.Minute / 15) // 15 requests per minute
	burst           = 5                            // allow bursts
	cleanupInterval = time.Minute * 5
)

func init() {
	go cleanupClients()
}

// Middleware to rate limit requests per IP
func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		limiter := getLimiter(ip)

		res := limiter.Reserve()
		resetTime := time.Now().Add(res.Delay())

		// Check allowance
		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", burst))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limit exceeded"}`))
			return
		}

		// Set headers for allowed requests
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", burst))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Burst()-1)) // approximate
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

		next.ServeHTTP(w, r)
	})
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if c, exists := clients[ip]; exists {
		c.lastSeen = time.Now()
		return c.limiter
	}

	// Correct: first param is rate.Limit, second is burst (int)
	lim := rate.NewLimiter(rateLimit, burst)
	clients[ip] = &client{
		limiter:  lim,
		lastSeen: time.Now(),
	}
	return lim
}

// Periodically clean old clients
func cleanupClients() {
	for {
		time.Sleep(cleanupInterval)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > cleanupInterval {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
