package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS returns a CORS middleware handler.
// allowedOrigins can be "*" for all origins or comma-separated list of specific origins.
func CORS(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is allowed
		allowAll := allowedOrigins == "*"
		originAllowed := false

		if allowAll {
			originAllowed = true
		} else {
			for _, o := range origins {
				if o == origin {
					originAllowed = true
					break
				}
			}
		}

		if originAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
