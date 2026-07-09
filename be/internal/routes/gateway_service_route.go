package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/handlers"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/middleware"
)

// RegisterGatewayServiceRoutes registers protected gateway service management routes.
func RegisterGatewayServiceRoutes(r *gin.RouterGroup, h *handlers.Handlers, acc *helpers.Access) {
	services := r.Group("/services")
	{
		services.POST("", middleware.RequirePermission(acc, "service.create"), h.GatewayServiceCreate)
		services.GET("", middleware.RequirePermission(acc, "service.index"), h.GatewayServiceGetAll)
		services.GET("/:id", middleware.RequirePermission(acc, "service.index"), h.GatewayServiceGetByID)
		services.PUT("/:id", middleware.RequirePermission(acc, "service.edit"), h.GatewayServiceUpdate)
		services.DELETE("/:id", middleware.RequirePermission(acc, "service.delete"), h.GatewayServiceDelete)
		services.POST("/:id/health-check", middleware.RequirePermission(acc, "service.health-check"), h.GatewayServiceHealthCheck)
	}
}
