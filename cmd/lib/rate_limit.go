// cmd/lib/rate_limit.go
package lib

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter structure pour limiter le nombre de requêtes
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]int
}

// NewRateLimiter crée un nouveau RateLimiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]int),
	}
}

// Limit vérifie et applique la limitation de débit
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		ip := r.RemoteAddr
		rl.visitors[ip]++

		if rl.visitors[ip] > 50 { // Limite à 10 requêtes
			http.Error(w, "Too many requests", http.StatusTooManyRequests) // a mettre avec le error.go HandleError
			return
		}

		go func() {
			time.Sleep(1 * time.Minute)
			rl.mu.Lock()
			rl.visitors[ip]--
			rl.mu.Unlock()
		}()

		next.ServeHTTP(w, r)
	})
}
