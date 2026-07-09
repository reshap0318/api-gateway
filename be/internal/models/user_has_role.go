package models

import "time"

// UserHasRole represents the many-to-many relationship between users and roles.
type UserHasRole struct {
	UserID     uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
	RoleID     uint      `gorm:"primaryKey;column:role_id" json:"role_id"`
	CreatedAt  time.Time `json:"created_at"`
	Role       Role      `gorm:"foreignKey:RoleID" json:"role"`
}

// TableName specifies the table name for UserHasRole model.
func (UserHasRole) TableName() string {
	return "user_has_roles"
}
