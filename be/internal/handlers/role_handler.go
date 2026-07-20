package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/dtos"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

// RoleCreate handles POST /api/roles
func (h *Handlers) RoleCreate(c *gin.Context) {
	var req dtos.RoleRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.RoleCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create role") {
		return
	}

	helpers.Created(c, "Role created successfully", dto)
}

// RoleGetAll handles GET /api/roles. page_size unset/negative (default) returns all
// records unpaginated; page_size>0 paginates. Response always includes metadata.
func (h *Handlers) RoleGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "-1"))

	if pageSize < 0 {
		page = 1
	}

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.RoleGetAll(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch roles") {
		return
	}

	helpers.OKWithMetadata(c, "Roles fetched successfully", result)
}

// RoleGetByID handles GET /api/roles/:id
func (h *Handlers) RoleGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	dto, err := h.svcs.RoleGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch role") {
		return
	}

	helpers.OK(c, "Role fetched successfully", dto)
}

// RoleUpdate handles PUT /api/roles/:id
func (h *Handlers) RoleUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	var req dtos.RoleRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.RoleUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update role") {
		return
	}

	helpers.OK(c, "Role updated successfully", dto)
}

// RoleDelete handles DELETE /api/roles/:id
func (h *Handlers) RoleDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	err = h.svcs.RoleDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete role") {
		return
	}

	helpers.OK(c, "Role deleted successfully", nil)
}

// RoleGetPermissions handles GET /api/roles/:id/permissions
func (h *Handlers) RoleGetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	perms, err := h.svcs.RoleGetPermissions(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch role permissions") {
		return
	}

	helpers.OK(c, "Role permissions fetched successfully", perms)
}
