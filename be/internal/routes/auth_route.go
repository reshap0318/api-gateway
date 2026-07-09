package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/handlers"
)

// RegisterAuthRoutes registers public authentication routes.
func RegisterAuthRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.AuthLogin)
		auth.POST("/refresh", h.AuthRefreshToken)
		auth.POST("/forgot-password", h.AuthForgetPassword)
		auth.POST("/reset-password", h.AuthResetPassword)
	}
}

// RegisterAuthProtectedRoutes registers protected authentication routes.
func RegisterAuthProtectedRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	auth := r.Group("/auth")
	{
		auth.POST("/logout", h.AuthLogout)
	}
}
