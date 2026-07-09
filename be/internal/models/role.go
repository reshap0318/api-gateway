package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a role in the system.
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description *string        `gorm:"size:255" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Permissions []Permission   `gorm:"many2many:role_has_permissions;joinForeignKey:role_id;joinReferences:permission_id" json:"permissions"`
}

// TableName specifies the table name for Role model.
func (Role) TableName() string {
	return "roles"
}
