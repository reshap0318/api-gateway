package repositories

import (
	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// RoleHasPermissionRepository provides database operations for RoleHasPermission model.
type RoleHasPermissionRepository struct {
	*GenericRepository[models.RoleHasPermission]
}

// NewRoleHasPermissionRepository creates a new RoleHasPermissionRepository.
func NewRoleHasPermissionRepository(db *gorm.DB) *RoleHasPermissionRepository {
	return &RoleHasPermissionRepository{
		GenericRepository: NewGenericRepository(db, &models.RoleHasPermission{}),
	}
}

// DeleteByRoleID deletes all role_has_permissions for a given role ID.
func (r *RoleHasPermissionRepository) DeleteByRoleID(tx *gorm.DB, roleID uint) error {
	return tx.Where("role_id = ?", roleID).Delete(&models.RoleHasPermission{}).Error
}
