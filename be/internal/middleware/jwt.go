package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/services"
)

// JWTAuth returns a JWT authentication middleware.
func JWTAuth(svcs *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := svcs.AuthValidateToken(token)
		if err != nil {
			helpers.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in request context and gin context
		ctx := context.WithValue(c.Request.Context(), helpers.KeyUserID, claims.UserID)
		c.Request = c.Request.WithContext(ctx)

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
