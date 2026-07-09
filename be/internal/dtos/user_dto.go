package dtos

import (
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/models"
)

// UserCreateRequest represents the request to create a user.
type UserCreateRequest struct {
	Name                 string `json:"name" validate:"required,min=2,max=100"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	Avatar               string `json:"avatar"`
	Roles                []uint `json:"roles"`
}

// UserUpdateRequest represents the request to update a user.
type UserUpdateRequest struct {
	Name                 string `json:"name" validate:"required,min=2,max=100"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"omitempty,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"omitempty,eqfield=Password"`
	Avatar               string `json:"avatar"`
	Roles                []uint `json:"roles"`
}

// UserStatusRequest represents the request to update a user's administrative status
// (FSD §2.25) — independent of the Lock mechanism (§2.26/§2.27).
type UserStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active suspended"`
}

// ToUserDTO converts User model to UserDTO.
func ToUserDTO(u *models.User) UserDTO {
	dto := UserDTO{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Avatar:      helpers.GetFileURL(u.Avatar),
		Status:      u.Status,
		LockedUntil: u.LockedUntil,
		CreatedAt:   u.CreatedAt,
		Roles:       []RoleMiniDTO{},
		Permissions: []PermissionDTO{},
	}

	permSet := make(map[uint]bool)
	for _, r := range u.Roles {
		dto.Roles = append(dto.Roles, ToRoleMiniDTO(&r))
		for _, p := range r.Permissions {
			if !permSet[p.ID] {
				permSet[p.ID] = true
				dto.Permissions = append(dto.Permissions, ToPermissionDTO(&p))
			}
		}
	}

	return dto
}

// ProfileUpdateRequest represents the request to update profile.
type ProfileUpdateRequest struct {
	Name                 string `json:"name" validate:"required,min=2,max=100"`
	Password             string `json:"password" validate:"omitempty,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"omitempty,eqfield=Password"`
	Avatar               string `json:"avatar"`
}

// ToUserDTOList converts a slice of User models to UserDTOs.
func ToUserDTOList(users []models.User) []UserDTO {
	result := make([]UserDTO, len(users))
	for i, u := range users {
		result[i] = ToUserDTO(&u)
	}
	return result
}
