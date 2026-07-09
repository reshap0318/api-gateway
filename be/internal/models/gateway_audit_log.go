package models

import "time"

// GatewayAuditLog represents an immutable audit trail entry for Service/Route config changes.
type GatewayAuditLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	EntityType  string    `gorm:"column:entity_type;size:20;not null" json:"entity_type"`
	EntityID    uint      `gorm:"column:entity_id;not null" json:"entity_id"`
	Action      string    `gorm:"size:20;not null" json:"action"`
	ActorUserID uint      `gorm:"column:actor_user_id;not null" json:"actor_user_id"`
	Changes     string    `gorm:"type:json" json:"changes"`
	CreatedAt   time.Time `json:"created_at"`

	Actor User `gorm:"foreignKey:ActorUserID" json:"-"`
}

// TableName specifies the table name for GatewayAuditLog model.
func (GatewayAuditLog) TableName() string {
	return "gateway_audit_logs"
}
