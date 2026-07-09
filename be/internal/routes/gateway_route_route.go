package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/handlers"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/middleware"
)

// RegisterGatewayRouteRoutes registers protected gateway route management routes.
func RegisterGatewayRouteRoutes(r *gin.RouterGroup, h *handlers.Handlers, acc *helpers.Access) {
	gwRoutes := r.Group("/routes")
	{
		gwRoutes.POST("", middleware.RequirePermission(acc, "route.create"), h.GatewayRouteCreate)
		gwRoutes.GET("", middleware.RequirePermission(acc, "route.index"), h.GatewayRouteGetAll)
		gwRoutes.GET("/:id", middleware.RequirePermission(acc, "route.index"), h.GatewayRouteGetByID)
		gwRoutes.PUT("/:id", middleware.RequirePermission(acc, "route.edit"), h.GatewayRouteUpdate)
		gwRoutes.DELETE("/:id", middleware.RequirePermission(acc, "route.delete"), h.GatewayRouteDelete)
	}
}

// RegisterGatewayCacheRoutes registers manual cache refresh/status routes.
// Permission: route.create OR route.edit (any-of) — no separate "gateway.cache-refresh"
// permission, per docs/CHANGELOG.md v1.1.0.
func RegisterGatewayCacheRoutes(r *gin.RouterGroup, h *handlers.Handlers, acc *helpers.Access) {
	cache := r.Group("/gateway/cache")
	{
		cache.POST("/refresh", middleware.RequirePermission(acc, "route.create", "route.edit"), h.GatewayCacheRefresh)
		cache.GET("/status", middleware.RequirePermission(acc, "route.create", "route.edit"), h.GatewayCacheStatus)
	}
}
