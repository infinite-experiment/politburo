package middleware

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var (
	limiters      = make(map[string]*rate.Limiter)
	limitersMutex sync.Mutex

	whitelistedIPs = map[string]bool{
		"127.0.0.1": true, // local bot
	}
)

func getLimiter(ip string) *rate.Limiter {
	limitersMutex.Lock()
	defer limitersMutex.Unlock()

	if limiter, exists := limiters[ip]; exists {
		return limiter
	}
	limiter := rate.NewLimiter(1, 5) // 1 request/sec, burst up to 5
	limiters[ip] = limiter
	return limiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if whitelistedIPs[ip] {
			next.ServeHTTP(w, r)
			return
		}

		limiter := getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
