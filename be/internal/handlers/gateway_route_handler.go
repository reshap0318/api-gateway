package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/dtos"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

// GatewayRouteCreate handles POST /api/routes
func (h *Handlers) GatewayRouteCreate(c *gin.Context) {
	var req dtos.GatewayRouteRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.GatewayRouteCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create route") {
		return
	}

	helpers.Created(c, "Route created successfully", result)
}

// GatewayRouteGetAll handles GET /api/routes
func (h *Handlers) GatewayRouteGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.GatewayRouteGetAllPaginated(
		c.Request.Context(),
		opts,
		c.Query("service"),
		c.Query("method"),
		c.Query("is_active"),
	)
	if helpers.HandleError(c, err, "Failed to fetch routes") {
		return
	}

	helpers.OKWithMetadata(c, "Routes fetched successfully", result)
}

// GatewayRouteGetByID handles GET /api/routes/:id
func (h *Handlers) GatewayRouteGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid route ID")
		return
	}

	result, err := h.svcs.GatewayRouteGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch route") {
		return
	}

	helpers.OK(c, "Route fetched successfully", result)
}

// GatewayRouteUpdate handles PUT /api/routes/:id
func (h *Handlers) GatewayRouteUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid route ID")
		return
	}

	var req dtos.GatewayRouteRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.GatewayRouteUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update route") {
		return
	}

	helpers.OK(c, "Route updated successfully", result)
}

// GatewayRouteDelete handles DELETE /api/routes/:id
func (h *Handlers) GatewayRouteDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid route ID")
		return
	}

	err = h.svcs.GatewayRouteDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete route") {
		return
	}

	helpers.OK(c, "Route deleted successfully", nil)
}
