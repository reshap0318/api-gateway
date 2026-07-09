package services

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/reshap0318/api-gateway/internal/dtos"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/models"
)

const healthCheckTimeout = 5 * time.Second

// GatewayServiceHealthCheck triggers an immediate health check for one Service (FSD §2.6).
// Not audit-logged and no notification is created — this is a status observation, not a
// configuration change.
func (s *Services) GatewayServiceHealthCheck(ctx context.Context, id uint) (*dtos.GatewayServiceHealthDTO, error) {
	service, err := s.repo.GatewayService.FindByID(nil, id)
	if err != nil {
		return nil, err
	}
	if !service.IsActive {
		return nil, &helpers.CustomError{Status: 400, Message: "Service tidak aktif"}
	}

	status := s.checkServiceHealth(service.BaseURL)
	now := time.Now()

	if _, err := s.repo.GatewayService.UpdateMap(nil, &models.GatewayService{ID: id}, map[string]interface{}{
		"health_status":     status,
		"health_checked_at": now,
	}); err != nil {
		return nil, err
	}

	return &dtos.GatewayServiceHealthDTO{HealthStatus: status, HealthCheckedAt: &now}, nil
}

// GatewayServiceHealthCheckAll runs the periodic background health check (FSD §2.22) against
// every active Service. A failure on one service is logged and does not abort the batch.
func (s *Services) GatewayServiceHealthCheckAll(ctx context.Context) {
	activeServices, err := s.repo.GatewayService.FindByFieldMap(nil, map[string]interface{}{"is_active": true})
	if err != nil {
		s.Logger.LogWarn("GatewayServiceHealthCheckAll", "Failed to list active services: %v", err)
		return
	}

	now := time.Now()
	for _, svc := range activeServices {
		status := s.checkServiceHealth(svc.BaseURL)
		if _, err := s.repo.GatewayService.UpdateMap(nil, &models.GatewayService{ID: svc.ID}, map[string]interface{}{
			"health_status":     status,
			"health_checked_at": now,
		}); err != nil {
			s.Logger.LogWarn("GatewayServiceHealthCheckAll", "Failed to update health for service %d: %v", svc.ID, err)
		}
	}

	s.Logger.LogInfo("GatewayServiceHealthCheckAll", "Health check completed for %d service(s)", len(activeServices))
}

// checkServiceHealth sends GET {base_url}/health with a short timeout and classifies the
// result: 2xx within the timeout → "up", anything else (including timeout) → "down".
func (s *Services) checkServiceHealth(baseURL string) string {
	client := http.Client{Timeout: healthCheckTimeout}

	resp, err := client.Get(strings.TrimRight(baseURL, "/") + "/health")
	if err != nil {
		return "down"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "up"
	}
	return "down"
}
