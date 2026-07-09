package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/handlers"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/middleware"
)

// RegisterPermissionRoutes registers protected permission routes.
func RegisterPermissionRoutes(r *gin.RouterGroup, h *handlers.Handlers, acc *helpers.Access) {
	permissions := r.Group("/permissions")
	{
		permissions.POST("", middleware.RequirePermission(acc, "permission.create"), h.PermissionCreate)
		permissions.GET("", middleware.RequirePermission(acc, "permission.index"), h.PermissionGetAll)
		permissions.GET("/:id", middleware.RequirePermission(acc, "permission.index"), h.PermissionGetByID)
		permissions.PUT("/:id", middleware.RequirePermission(acc, "permission.edit"), h.PermissionUpdate)
		permissions.DELETE("/:id", middleware.RequirePermission(acc, "permission.delete"), h.PermissionDelete)
	}
}
