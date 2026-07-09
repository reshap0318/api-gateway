package dtos

import (
	"time"

	"github.com/reshap0318/api-gateway/internal/models"
)

// GatewayRouteRequest represents the request to create or update a gateway route.
// Field naming convention: no "_id"/"_ids" suffix on request payloads (see docs/CHANGELOG.md v1.1.0).
type GatewayRouteRequest struct {
	Service             uint   `json:"service" validate:"required"`
	Method              string `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE *"`
	PathPattern         string `json:"path_pattern" validate:"required,min=1"`
	PermissionMatchMode string `json:"permission_match_mode" validate:"omitempty,oneof=any all"`
	Permissions         []uint `json:"permissions"`
	RateLimitPerMinute  *int   `json:"rate_limit_per_minute" validate:"omitempty,min=1"`
	IsActive            *bool  `json:"is_active"`
}

// GatewayRouteDTO represents gateway route data transfer object.
type GatewayRouteDTO struct {
	ID                  uint                  `json:"id"`
	ServiceID           uint                  `json:"service_id"`
	Service             GatewayServiceMiniDTO `json:"service"`
	Method              string                `json:"method"`
	PathPattern         string                `json:"path_pattern"`
	PermissionMatchMode string                `json:"permission_match_mode"`
	Permissions         []PermissionDTO       `json:"permissions"`
	RateLimitPerMinute  *int                  `json:"rate_limit_per_minute"`
	IsActive            bool                  `json:"is_active"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}

// ToGatewayRouteDTO converts GatewayRoute model to GatewayRouteDTO.
func ToGatewayRouteDTO(r *models.GatewayRoute) GatewayRouteDTO {
	dto := GatewayRouteDTO{
		ID:                  r.ID,
		ServiceID:           r.ServiceID,
		Method:              r.Method,
		PathPattern:         r.PathPattern,
		PermissionMatchMode: r.PermissionMatchMode,
		Permissions:         []PermissionDTO{},
		RateLimitPerMinute:  r.RateLimitPerMinute,
		IsActive:            r.IsActive,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
	}

	if r.Service.ID != 0 {
		dto.Service = ToGatewayServiceMiniDTO(&r.Service)
	}

	for _, p := range r.Permissions {
		dto.Permissions = append(dto.Permissions, ToPermissionDTO(&p))
	}

	return dto
}

// ToGatewayRouteDTOList converts a slice of GatewayRoute models to DTOs.
func ToGatewayRouteDTOList(routes []models.GatewayRoute) []GatewayRouteDTO {
	result := make([]GatewayRouteDTO, len(routes))
	for i, r := range routes {
		result[i] = ToGatewayRouteDTO(&r)
	}
	return result
}
