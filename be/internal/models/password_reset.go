package models

import (
	"time"

	"gorm.io/gorm"
)

// PasswordReset represents a password reset token.
type PasswordReset struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"index;size:255;not null" json:"email"`
	Token     string         `gorm:"uniqueIndex;size:255;not null" json:"-"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	Used      bool           `gorm:"default:false" json:"used"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for PasswordReset model.
func (PasswordReset) TableName() string {
	return "password_resets"
}
