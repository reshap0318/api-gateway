package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/helpers"
)

// RequirePermission returns middleware that checks if user has ANY of the specified permissions.
func RequirePermission(acc *helpers.Access, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !acc.HasPermission(c.Request.Context(), permissions...) {
			helpers.Forbidden(c, "Permission denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireRole returns middleware that checks if user has the specified role.
func RequireRole(acc *helpers.Access, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !acc.HasRole(c.Request.Context(), role) {
			helpers.Forbidden(c, "Role required")
			c.Abort()
			return
		}
		c.Next()
	}
}
