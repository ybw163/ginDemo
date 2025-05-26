package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

var (
	limiter = rate.NewLimiter(100, 200) // 每秒100个请求，突发100个
	mu      sync.RWMutex
	clients = make(map[string]*rate.Limiter)
)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		mu.RLock()
		lmt, exists := clients[clientIP]
		mu.RUnlock()

		if !exists {
			mu.Lock()
			clients[clientIP] = limiter
			mu.Unlock()
		}

		if !lmt.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
