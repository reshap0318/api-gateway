package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

func (h *Handlers) NotificationGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		Order:    "DESC",
	}

	var isRead *bool
	if isReadStr := c.Query("is_read"); isReadStr != "" {
		if isReadStr == "true" {
			val := true
			isRead = &val
		} else if isReadStr == "false" {
			val := false
			isRead = &val
		}
	}

	notifType := c.Query("type")

	if isRead != nil || notifType != "" {
		result, err := h.svcs.NotificationGetAllWithFilters(c.Request.Context(), opts, isRead, notifType)
		if helpers.HandleError(c, err, "Failed to fetch notifications") {
			return
		}

		helpers.OKWithMetadata(c, "Notifications fetched successfully", result)
		return
	}

	result, err := h.svcs.NotificationGetAll(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch notifications") {
		return
	}

	helpers.OKWithMetadata(c, "Notifications fetched successfully", result)
}

func (h *Handlers) NotificationGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid notification ID")
		return
	}

	dto, err := h.svcs.NotificationGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch notification") {
		return
	}

	helpers.OK(c, "Notification fetched successfully", dto)
}

func (h *Handlers) NotificationMarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid notification ID")
		return
	}

	err = h.svcs.NotificationMarkAsRead(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to mark notification as read") {
		return
	}

	helpers.OK(c, "Notification marked as read", nil)
}

func (h *Handlers) NotificationMarkAllAsRead(c *gin.Context) {
	err := h.svcs.NotificationMarkAllAsRead(c.Request.Context())
	if helpers.HandleError(c, err, "Failed to mark all notifications as read") {
		return
	}

	helpers.OK(c, "All notifications marked as read", nil)
}

func (h *Handlers) NotificationCountUnread(c *gin.Context) {
	count, err := h.svcs.NotificationCountUnread(c.Request.Context())
	if helpers.HandleError(c, err, "Failed to count unread notifications") {
		return
	}

	helpers.OK(c, "Unread count fetched successfully", map[string]int64{
		"unread_count": count,
	})
}

func (h *Handlers) NotificationDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid notification ID")
		return
	}

	err = h.svcs.NotificationDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete notification") {
		return
	}

	helpers.OK(c, "Notification deleted successfully", nil)
}
