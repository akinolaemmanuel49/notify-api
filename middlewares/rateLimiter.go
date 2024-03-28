package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/akinolaemmanuel49/notify-api/config"
	"github.com/akinolaemmanuel49/notify-api/utils"
)

var (
	requestCount      = make(map[string]int)
	mu                sync.Mutex
	last429Reset      = time.Now()
	durationSinceLast = 0 * time.Second
	rateLimitExceeded = false
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	utils.LoadEnv()
	var cfg config.Config

	cfg.ReadFile("dev-config.yml") // For use in development
	cfg.ReadEnv()

	maxRequests, err := strconv.Atoi(cfg.RateLimiting.MaxRequests)
	if err != nil {
		// Handle error
		maxRequests = 100 // Set a default value
	}

	durationConv, err := strconv.Atoi(cfg.RateLimiting.Duration)
	if err != nil {
		// Handle error
		durationConv = 5 // Set a default value (assuming duration is in minutes)
	}
	duration := time.Minute * time.Duration(durationConv)

	// Return the middleware handler.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the IP address from RemoteAddr (IPv4 or IPv6).
		ip := strings.Split(r.RemoteAddr, ":")[0]

		// Reset request count if the rate limit was exceeded before and duration has passed since last 429
		if rateLimitExceeded && time.Since(last429Reset) >= durationSinceLast {
			mu.Lock()
			requestCount = make(map[string]int)
			rateLimitExceeded = false
			mu.Unlock()
		}

		// Increment the request count for this IP in a thread-safe manner.
		mu.Lock()
		count := requestCount[ip]
		requestCount[ip]++
		mu.Unlock()

		// Check if the request count exceeds the maximum allowed.
		if count > maxRequests {
			utils.RespondWithError(w, "Rate limit exceeded", http.StatusTooManyRequests)
			if !rateLimitExceeded {
				rateLimitExceeded = true
				last429Reset = time.Now()
				durationSinceLast = duration
			}
			return
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
