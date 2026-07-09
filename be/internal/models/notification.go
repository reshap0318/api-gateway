package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Type      string         `gorm:"size:100;not null" json:"type"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Data      string         `gorm:"type:json" json:"data"`
	ReadAt    *time.Time     `json:"read_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (Notification) TableName() string {
	return "notifications"
}
