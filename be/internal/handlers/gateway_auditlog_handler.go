package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

// GatewayAuditLogGetAll handles GET /api/audit-logs
func (h *Handlers) GatewayAuditLogGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.GatewayAuditLogGetAllPaginated(
		c.Request.Context(),
		opts,
		c.Query("entity_type"),
		c.Query("entity"),
		c.Query("actor"),
		c.Query("from"),
		c.Query("to"),
	)
	if helpers.HandleError(c, err, "Failed to fetch audit logs") {
		return
	}

	helpers.OKWithMetadata(c, "Audit logs fetched successfully", result)
}

// GatewayAuditLogGetByID handles GET /api/audit-logs/:id
func (h *Handlers) GatewayAuditLogGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid audit log ID")
		return
	}

	result, err := h.svcs.GatewayAuditLogGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch audit log") {
		return
	}

	helpers.OK(c, "Audit log fetched successfully", result)
}
