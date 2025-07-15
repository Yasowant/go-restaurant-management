package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	tokens   int
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func RateLimitMiddleware() gin.HandlerFunc {
	go cleanupVisitors() // Start background cleanup

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			v = &visitor{lastSeen: time.Now(), tokens: 10}
			visitors[ip] = v
		}

		// refill token if 1 second has passed
		if time.Since(v.lastSeen) > time.Second {
			v.tokens = 10
			v.lastSeen = time.Now()
		}

		if v.tokens <= 0 {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		v.tokens--
		mu.Unlock()

		c.Next()
	}
}

func cleanupVisitors() {
	for {
		time.Sleep(10 * time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 15*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
