package models

import (
	"time"

	"gorm.io/gorm"
)

// GatewayService represents an upstream service registered for reverse proxying.
type GatewayService struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	BaseURL            string         `gorm:"column:base_url;size:500;not null" json:"base_url"`
	BasePath           string         `gorm:"column:base_path;size:200;not null;uniqueIndex" json:"base_path"`
	Protocol           string         `gorm:"size:20;not null;default:http" json:"protocol"`
	IsActive           bool           `gorm:"column:is_active;not null;default:true" json:"is_active"`
	RateLimitPerMinute *int           `gorm:"column:rate_limit_per_minute" json:"rate_limit_per_minute"`
	HealthStatus       string         `gorm:"column:health_status;size:20;not null;default:unknown" json:"health_status"`
	HealthCheckedAt    *time.Time     `gorm:"column:health_checked_at" json:"health_checked_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	Routes             []GatewayRoute `gorm:"foreignKey:ServiceID" json:"routes,omitempty"`
}

// TableName specifies the table name for GatewayService model.
func (GatewayService) TableName() string {
	return "gateway_services"
}
