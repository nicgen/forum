package RateLimit

import (
	"forum/handlers"
	"forum/models"
	"net/http"
	"sync"
	"time"
)

// RateLimiter structure pour limiter le nombre de requêtes
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]map[string]int
}

// NewRateLimiter crée un nouveau RateLimiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]map[string]int),
	}
}

// Limit vérifie et applique la limitation de débit
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		ip := r.RemoteAddr
		path := r.URL.Path

		if rl.visitors[ip] == nil {
			rl.visitors[ip] = make(map[string]int)
		}
		rl.visitors[ip][path]++

		if rl.visitors[ip][path] > 100 { // Limite à 100 requêtes
			// Erreur personnalisé via error.go
			err := &models.CustomError{
				StatusCode: http.StatusTooManyRequests,
				Message:    "Too Many Requests",
			}
			//Appel de HandleError
			handlers.HandleError(w, err.StatusCode, err.Message)
			return
		}

		go func() {
			time.Sleep(1 * time.Minute)
			rl.mu.Lock()
			rl.visitors[ip][path]--
			if rl.visitors[ip][path] <= 0 {
				delete(rl.visitors[ip], path)
				if len(rl.visitors[ip]) == 0 {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}()

		next.ServeHTTP(w, r)
	})
}
