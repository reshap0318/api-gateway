package dtos

import (
	"time"

	"github.com/reshap0318/api-gateway/internal/models"
)

// GatewayAuditLogDTO represents an audit trail entry data transfer object.
type GatewayAuditLogDTO struct {
	ID          uint      `json:"id"`
	EntityType  string    `json:"entity_type"`
	EntityID    uint      `json:"entity_id"`
	Action      string    `json:"action"`
	ActorUserID uint      `json:"actor_user_id"`
	ActorName   string    `json:"actor_name"`
	Changes     string    `json:"changes"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToGatewayAuditLogDTO converts GatewayAuditLog model to GatewayAuditLogDTO.
func ToGatewayAuditLogDTO(a *models.GatewayAuditLog) GatewayAuditLogDTO {
	return GatewayAuditLogDTO{
		ID:          a.ID,
		EntityType:  a.EntityType,
		EntityID:    a.EntityID,
		Action:      a.Action,
		ActorUserID: a.ActorUserID,
		ActorName:   a.Actor.Name,
		Changes:     a.Changes,
		CreatedAt:   a.CreatedAt,
	}
}

// ToGatewayAuditLogDTOList converts a slice of GatewayAuditLog models to DTOs.
func ToGatewayAuditLogDTOList(logs []models.GatewayAuditLog) []GatewayAuditLogDTO {
	result := make([]GatewayAuditLogDTO, len(logs))
	for i, a := range logs {
		result[i] = ToGatewayAuditLogDTO(&a)
	}
	return result
}
