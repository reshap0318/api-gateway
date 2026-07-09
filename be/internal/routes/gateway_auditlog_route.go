package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/handlers"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/middleware"
)

// RegisterGatewayAuditLogRoutes registers protected audit trail read-only routes.
func RegisterGatewayAuditLogRoutes(r *gin.RouterGroup, h *handlers.Handlers, acc *helpers.Access) {
	auditLogs := r.Group("/audit-logs")
	{
		auditLogs.GET("", middleware.RequirePermission(acc, "audit.index"), h.GatewayAuditLogGetAll)
		auditLogs.GET("/:id", middleware.RequirePermission(acc, "audit.index"), h.GatewayAuditLogGetByID)
	}
}
