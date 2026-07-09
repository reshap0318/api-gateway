package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reshap0318/api-gateway/internal/helpers"
)

// HealthCheck handles health check endpoint.
// @Summary Health check
// @Description Check system health status (database, redis, etc)
// @Tags health
// @Produce json
// @Success 200 {object} dtos.HealthStatus
// @Router /api/health [get]
func (h *Handlers) HealthCheck(c *gin.Context) {
	status := h.svcs.HealthGetStatus(c.Request.Context())

	if status.Status == "unhealthy" {
		helpers.ErrorResponse(c, http.StatusServiceUnavailable, "Service unavailable")
		return
	}

	helpers.OK(c, "Service healthy", status)
}
