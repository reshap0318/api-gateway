package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Email               string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password            string         `gorm:"size:255;not null" json:"-"`
	Name                string         `gorm:"size:255" json:"name"`
	Avatar              string         `gorm:"size:500" json:"avatar"`
	Status              string         `gorm:"size:20;not null;default:active" json:"status"`
	FailedLoginAttempts int            `gorm:"column:failed_login_attempts;not null;default:0" json:"-"`
	LockedUntil         *time.Time     `gorm:"column:locked_until" json:"-"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
	Roles               []Role         `gorm:"many2many:user_has_roles;joinForeignKey:user_id;joinReferences:role_id" json:"roles"`
}

// TableName specifies the table name for User model.
func (User) TableName() string {
	return "users"
}
