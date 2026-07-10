package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/helpers"
)

// RateLimit returns a rate limiting middleware that applies globally.
// It sets rate limit headers on all responses and returns 429 when limit exceeded.
func RateLimit(limiter *helpers.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.ClientIP() resolves the real client IP from X-Forwarded-For/X-Real-IP only when
		// the immediate peer is in gin's configured TrustedProxies (see main.go) — trusting
		// those headers unconditionally would let any client spoof a fresh IP per request and
		// bypass rate limiting entirely.
		ip := c.ClientIP()

		// Check rate limit
		info := limiter.Allow(ip)

		// Set rate limit headers on ALL responses
		c.Header("X-RateLimit-Remaining", strconv.Itoa(info.Remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(info.Reset, 10))

		if !info.Allowed {
			// Calculate retry-after in seconds
			retryAfter := int(info.Reset - time.Now().Unix())
			if retryAfter < 0 {
				retryAfter = 0
			}
			c.Header("Retry-After", strconv.Itoa(retryAfter))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
