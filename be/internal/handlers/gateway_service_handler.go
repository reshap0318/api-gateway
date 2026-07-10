package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/dtos"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

// GatewayServiceCreate handles POST /api/services
func (h *Handlers) GatewayServiceCreate(c *gin.Context) {
	var req dtos.GatewayServiceRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.GatewayServiceCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create service") {
		return
	}

	helpers.Created(c, "Service created successfully", result)
}

// GatewayServiceGetAll handles GET /api/services
func (h *Handlers) GatewayServiceGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "0"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.GatewayServiceGetAllPaginated(
		c.Request.Context(),
		opts,
		c.Query("search"),
		c.Query("protocol"),
		c.Query("is_active"),
		c.Query("health_status"),
	)
	if helpers.HandleError(c, err, "Failed to fetch services") {
		return
	}

	helpers.OKWithMetadata(c, "Services fetched successfully", result)
}

// GatewayServiceGetByID handles GET /api/services/:id
func (h *Handlers) GatewayServiceGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid service ID")
		return
	}

	result, err := h.svcs.GatewayServiceGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch service") {
		return
	}

	helpers.OK(c, "Service fetched successfully", result)
}

// GatewayServiceUpdate handles PUT /api/services/:id
func (h *Handlers) GatewayServiceUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid service ID")
		return
	}

	var req dtos.GatewayServiceRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.GatewayServiceUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update service") {
		return
	}

	helpers.OK(c, "Service updated successfully", result)
}

// GatewayServiceHealthCheck handles POST /api/services/:id/health-check
func (h *Handlers) GatewayServiceHealthCheck(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid service ID")
		return
	}

	result, err := h.svcs.GatewayServiceHealthCheck(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to run health check") {
		return
	}

	helpers.OK(c, "Health check completed", result)
}

// GatewayServiceDelete handles DELETE /api/services/:id
func (h *Handlers) GatewayServiceDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid service ID")
		return
	}

	cascade := c.Query("cascade") == "true"

	err = h.svcs.GatewayServiceDelete(c.Request.Context(), uint(id), cascade)
	if helpers.HandleError(c, err, "Failed to delete service") {
		return
	}

	helpers.OK(c, "Service deleted successfully", nil)
}
