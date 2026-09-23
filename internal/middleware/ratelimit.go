package middleware

import (
	"net/http"
	"sync"
	"time"

	"movie-ticket/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	sync.RWMutex
	ips   map[string]*ipLimiter
	rate  rate.Limit
	burst int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	limiter := &RateLimiter{
		ips:   make(map[string]*ipLimiter),
		rate:  r,
		burst: b,
	}

	// Clean up idle IPs periodically
	go func() {
		for {
			time.Sleep(3 * time.Minute)
			limiter.Lock()
			for ip, client := range limiter.ips {
				if time.Since(client.lastSeen) > 5*time.Minute {
					delete(limiter.ips, ip)
				}
			}
			limiter.Unlock()
		}
	}()

	return limiter
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.Lock()
	defer rl.Unlock()

	lim, exists := rl.ips[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.ips[ip] = &ipLimiter{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	lim.lastSeen = time.Now()
	return lim.limiter
}

func RateLimitMiddleware(r rate.Limit, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(r, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, models.StandardResponse{
				Success: false,
				Message: "Too many requests. Please slow down.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
