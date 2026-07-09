package dtos

import "github.com/reshap0318/api-gateway/internal/models"

// RoleRequest represents the request to create or update a role.
type RoleRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description *string `json:"description" validate:"omitempty,max=255"`
	Permissions []uint  `json:"permissions" validate:"required"`
}

// RoleDTO represents role data transfer object.
type RoleDTO struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Permissions []PermissionDTO `json:"permissions"`
}

// ToRoleDTO converts Role model to RoleDTO.
func ToRoleDTO(r *models.Role) RoleDTO {
	dto := RoleDTO{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: []PermissionDTO{},
	}

	for _, p := range r.Permissions {
		dto.Permissions = append(dto.Permissions, ToPermissionDTO(&p))
	}

	return dto
}

// ToRoleDTOList converts a slice of Role models to RoleDTOs.
func ToRoleDTOList(roles []models.Role) []RoleDTO {
	result := make([]RoleDTO, len(roles))
	for i, r := range roles {
		result[i] = ToRoleDTO(&r)
	}
	return result
}

// RoleMiniDTO represents a lightweight role for user response.
type RoleMiniDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// ToRoleMiniDTO converts Role model to RoleMiniDTO.
func ToRoleMiniDTO(r *models.Role) RoleMiniDTO {
	return RoleMiniDTO{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
}
