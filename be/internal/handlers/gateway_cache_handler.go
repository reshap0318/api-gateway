package handlers

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/helpers"
)

// GatewayCacheRefresh handles POST /api/gateway/cache/refresh (permission: route.create OR route.edit).
func (h *Handlers) GatewayCacheRefresh(c *gin.Context) {
	if h.RouteManager == nil {
		helpers.InternalServerError(c, "Route manager not initialized")
		return
	}

	if err := h.RouteManager.Refresh(); err != nil {
		helpers.InternalServerError(c, "Failed to refresh cache")
		return
	}

	helpers.OK(c, "Cache refreshed", gin.H{"refreshed_at": time.Now()})
}

// GatewayCacheStatus handles GET /api/gateway/cache/status (permission: route.create OR route.edit).
func (h *Handlers) GatewayCacheStatus(c *gin.Context) {
	if h.RouteManager == nil {
		helpers.InternalServerError(c, "Route manager not initialized")
		return
	}

	totalRoutes, totalServices, lastRefreshed, instanceID := h.RouteManager.Stats()

	helpers.OK(c, "Cache status fetched", gin.H{
		"last_refreshed_at": lastRefreshed,
		"total_services":    totalServices,
		"total_routes":      totalRoutes,
		"instance_id":       instanceID,
	})
}
