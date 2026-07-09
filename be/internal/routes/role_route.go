package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/handlers"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/middleware"
)

// RegisterRoleRoutes registers protected role routes.
func RegisterRoleRoutes(r *gin.RouterGroup, h *handlers.Handlers, acc *helpers.Access) {
	roles := r.Group("/roles")
	{
		roles.POST("", middleware.RequirePermission(acc, "role.create"), h.RoleCreate)
		roles.GET("", middleware.RequirePermission(acc, "role.index"), h.RoleGetAll)
		roles.GET("/:id", middleware.RequirePermission(acc, "role.index"), h.RoleGetByID)
		roles.PUT("/:id", middleware.RequirePermission(acc, "role.edit"), h.RoleUpdate)
		roles.DELETE("/:id", middleware.RequirePermission(acc, "role.delete"), h.RoleDelete)
		roles.GET("/:id/permissions", middleware.RequirePermission(acc, "role.index"), h.RoleGetPermissions)
	}
}
