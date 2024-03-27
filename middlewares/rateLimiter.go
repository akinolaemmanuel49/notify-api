package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/akinolaemmanuel49/notify-api/config"
	"github.com/akinolaemmanuel49/notify-api/utils"
)

var (
	requestCount = make(map[string]int)
	mu           sync.Mutex
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	utils.LoadEnv()
	var cfg config.Config

	cfg.ReadFile("dev-config.yml") // For use in development
	cfg.ReadEnv()

	maxRequests, _ := strconv.Atoi(cfg.RateLimiting.MaxRequests)
	durationConv, _ := strconv.Atoi(cfg.RateLimiting.Duration)
	duration := time.Minute * time.Duration(durationConv)

	// Return the middleware handler.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the IP address from RemoteAddr (IPv4 or IPv6).
		ip := strings.Split(r.RemoteAddr, ":")[0]

		// Increment the request count for this IP in a thread-safe manner.
		mu.Lock()
		requestCount[ip]++
		count := requestCount[ip]
		fmt.Println(requestCount)
		mu.Unlock()

		// Reset the request count after the specified duration.
		time.AfterFunc(duration, func() {
			mu.Lock()
			defer mu.Unlock()
			requestCount[ip] = 0
		})

		// Check if the request count exceeds the maximum allowed.
		if count > maxRequests {
			utils.RespondWithError(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
