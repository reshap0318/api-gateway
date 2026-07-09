package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/api-gateway/internal/dtos"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

// PermissionCreate handles POST /api/permissions
func (h *Handlers) PermissionCreate(c *gin.Context) {
	var req dtos.PermissionRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.PermissionCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create permission") {
		return
	}

	helpers.Created(c, "Permission created successfully", dto)
}

// PermissionGetAll handles GET /api/permissions with optional pagination
func (h *Handlers) PermissionGetAll(c *gin.Context) {
	pageStr := c.Query("page")

	if pageStr == "" {
		permissions, err := h.svcs.PermissionGetAll(c.Request.Context())
		if helpers.HandleError(c, err, "Failed to fetch permissions") {
			return
		}

		helpers.OK(c, "Permissions fetched successfully", permissions)
		return
	}

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.PermissionGetAllPaginated(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch permissions") {
		return
	}

	helpers.OKWithMetadata(c, "Permissions fetched successfully", result)
}

// PermissionGetByID handles GET /api/permissions/:id
func (h *Handlers) PermissionGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	dto, err := h.svcs.PermissionGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch permission") {
		return
	}

	helpers.OK(c, "Permission fetched successfully", dto)
}

// PermissionUpdate handles PUT /api/permissions/:id
func (h *Handlers) PermissionUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	var req dtos.PermissionRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	dto, err := h.svcs.PermissionUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update permission") {
		return
	}

	helpers.OK(c, "Permission updated successfully", dto)
}

// PermissionDelete handles DELETE /api/permissions/:id
func (h *Handlers) PermissionDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	err = h.svcs.PermissionDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete permission") {
		return
	}

	helpers.OK(c, "Permission deleted successfully", nil)
}
