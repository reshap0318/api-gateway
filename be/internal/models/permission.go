package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission represents a permission in the system.
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description *string        `gorm:"size:255" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Permission model.
func (Permission) TableName() string {
	return "permissions"
}
