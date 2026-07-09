package models

import (
	"time"

	"gorm.io/gorm"
)

// GatewayRoute represents a routing rule (method + path pattern) under a GatewayService.
type GatewayRoute struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	ServiceID           uint           `gorm:"column:service_id;not null" json:"service_id"`
	Method              string         `gorm:"size:10;not null" json:"method"`
	PathPattern         string         `gorm:"column:path_pattern;size:500;not null" json:"path_pattern"`
	PermissionMatchMode string         `gorm:"column:permission_match_mode;size:10;not null;default:any" json:"permission_match_mode"`
	RateLimitPerMinute  *int           `gorm:"column:rate_limit_per_minute" json:"rate_limit_per_minute"`
	IsActive            bool           `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
	Service             GatewayService `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	Permissions         []Permission   `gorm:"many2many:gateway_route_permissions;joinForeignKey:route_id;joinReferences:permission_id" json:"permissions"`
}

// TableName specifies the table name for GatewayRoute model.
func (GatewayRoute) TableName() string {
	return "gateway_routes"
}
