package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/handlers"
)

// RegisterSystemRoutes registers system-level public routes (health, jwks).
func RegisterSystemRoutes(r *gin.Engine, h *handlers.Handlers) {
	r.GET("/health", h.HealthCheck)
	r.GET("/.well-known/jwks.json", h.JWKSGetKeys)
}

// RegisterSystemProtectedRoutes registers system-level protected routes.
func RegisterSystemProtectedRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	r.POST("/upload", h.FileUpload)
}
