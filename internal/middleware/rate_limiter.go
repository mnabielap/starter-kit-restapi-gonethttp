package middleware

import (
	"net/http"
	"sync"

	"starter-kit-restapi-gonethttp/pkg/response"

	"golang.org/x/time/rate"
)

// IPRateLimiter holds rate limiters for each IP address
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter creates a limiter that allows r events per second with burst b
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// GetLimiter returns the rate limiter for the provided IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// RateLimit middleware limits requests based on IP
func RateLimit(next http.Handler) http.Handler {
	// Allow 5 requests per second with burst of 20
	limiter := NewIPRateLimiter(5, 20)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if !limiter.GetLimiter(ip).Allow() {
			response.Error(w, http.StatusTooManyRequests, "Too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}