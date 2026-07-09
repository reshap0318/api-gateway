package models

import (
	"time"
)

// RoleHasPermission represents the many-to-many relationship between roles and permissions.
type RoleHasPermission struct {
	RoleID       uint        `gorm:"primaryKey;column:role_id" json:"role_id"`
	PermissionID uint        `gorm:"primaryKey;column:permission_id" json:"permission_id"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Role         Role        `gorm:"foreignKey:RoleID" json:"-"`
	Permission   Permission  `gorm:"foreignKey:PermissionID" json:"-"`
}

// TableName specifies the table name for RoleHasPermission model.
func (RoleHasPermission) TableName() string {
	return "role_has_permissions"
}
