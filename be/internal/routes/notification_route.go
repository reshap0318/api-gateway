package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/handlers"
)

func RegisterNotificationRoutes(r *gin.RouterGroup, h *handlers.Handlers) {
	notifications := r.Group("/notifications")
	{
		notifications.GET("", h.NotificationGetAll)
		notifications.GET("/unread-count", h.NotificationCountUnread)
		notifications.GET("/:id", h.NotificationGetByID)
		notifications.PATCH("/:id/read", h.NotificationMarkAsRead)
		notifications.PATCH("/mark-all-read", h.NotificationMarkAllAsRead)
		notifications.DELETE("/:id", h.NotificationDelete)
	}
}
