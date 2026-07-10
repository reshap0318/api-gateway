package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/services"
)

// ExtractBearerClaims reads the "Authorization: Bearer <token>" header and validates the JWT.
// Shared by JWTAuth (Management API, scoped to the `/api` group) and the Dynamic Proxy
// Engine's proxy.authenticate (registered on gin's NoRoute, which bypasses group-scoped
// middleware entirely) so the two entry points can't drift on what counts as "authenticated" —
// each caller still decides what to do with the resulting claims afterwards.
func ExtractBearerClaims(c *gin.Context, svcs *services.Services) (*helpers.JWTClaims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("Authorization header required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("Invalid authorization header format")
	}

	claims, err := svcs.AuthValidateToken(parts[1])
	if err != nil {
		return nil, errors.New("Invalid or expired token")
	}

	return claims, nil
}

// JWTAuth returns a JWT authentication middleware.
func JWTAuth(svcs *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ExtractBearerClaims(c, svcs)
		if err != nil {
			helpers.Unauthorized(c, err.Error())
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
