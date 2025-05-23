package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

var (
	limiter = rate.NewLimiter(10, 100) // 每秒10个请求，突发100个
	mu      sync.RWMutex
	clients = make(map[string]*rate.Limiter)
)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		mu.RLock()
		limiter, exists := clients[clientIP]
		mu.RUnlock()

		if !exists {
			mu.Lock()
			limiter = rate.NewLimiter(10, 100)
			clients[clientIP] = limiter
			mu.Unlock()
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
