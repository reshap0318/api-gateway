package dtos

import (
	"time"

	"github.com/reshap0318/api-gateway/internal/models"
)

// GatewayServiceRequest represents the request to create or update a gateway service.
type GatewayServiceRequest struct {
	Name               string `json:"name" validate:"required,min=3,max=100"`
	BaseURL            string `json:"base_url" validate:"required,url"`
	Protocol           string `json:"protocol" validate:"required,oneof=http websocket"`
	RateLimitPerMinute *int   `json:"rate_limit_per_minute" validate:"omitempty,min=1"`
	IsActive           *bool  `json:"is_active"`
}

// GatewayServiceDTO represents gateway service data transfer object.
type GatewayServiceDTO struct {
	ID                 uint       `json:"id"`
	Name               string     `json:"name"`
	BaseURL            string     `json:"base_url"`
	Protocol           string     `json:"protocol"`
	IsActive           bool       `json:"is_active"`
	RateLimitPerMinute *int       `json:"rate_limit_per_minute"`
	HealthStatus       string     `json:"health_status"`
	HealthCheckedAt    *time.Time `json:"health_checked_at"`
	RouteCount         int        `json:"route_count"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// ToGatewayServiceDTO converts GatewayService model to GatewayServiceDTO.
func ToGatewayServiceDTO(s *models.GatewayService) GatewayServiceDTO {
	return GatewayServiceDTO{
		ID:                 s.ID,
		Name:               s.Name,
		BaseURL:            s.BaseURL,
		Protocol:           s.Protocol,
		IsActive:           s.IsActive,
		RateLimitPerMinute: s.RateLimitPerMinute,
		HealthStatus:       s.HealthStatus,
		HealthCheckedAt:    s.HealthCheckedAt,
		RouteCount:         len(s.Routes),
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
	}
}

// ToGatewayServiceDTOList converts a slice of GatewayService models to DTOs.
func ToGatewayServiceDTOList(services []models.GatewayService) []GatewayServiceDTO {
	result := make([]GatewayServiceDTO, len(services))
	for i, s := range services {
		result[i] = ToGatewayServiceDTO(&s)
	}
	return result
}

// GatewayServiceHealthDTO represents the response of a health check trigger (TDD §4.2).
type GatewayServiceHealthDTO struct {
	HealthStatus    string     `json:"health_status"`
	HealthCheckedAt *time.Time `json:"health_checked_at"`
}

// GatewayServiceMiniDTO represents a lightweight service reference (used in RouteDTO).
type GatewayServiceMiniDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

// ToGatewayServiceMiniDTO converts GatewayService model to GatewayServiceMiniDTO.
func ToGatewayServiceMiniDTO(s *models.GatewayService) GatewayServiceMiniDTO {
	return GatewayServiceMiniDTO{
		ID:       s.ID,
		Name:     s.Name,
		IsActive: s.IsActive,
	}
}
