package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Simple in-memory token bucket rate limiter per IP.
type bucket struct {
	tokens   float64
	lastTime time.Time
}

func RateLimit(rps float64, burst int) gin.HandlerFunc {
	var mu sync.Mutex
	buckets := make(map[string]*bucket)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		b, ok := buckets[ip]
		if !ok {
			b = &bucket{tokens: float64(burst), lastTime: now}
			buckets[ip] = b
		}

		elapsed := now.Sub(b.lastTime).Seconds()
		b.tokens += elapsed * rps
		if b.tokens > float64(burst) {
			b.tokens = float64(burst)
		}
		b.lastTime = now

		if b.tokens < 1 {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "too many requests, please slow down",
				},
			})
			return
		}
		b.tokens--
		mu.Unlock()

		c.Next()
	}
}
