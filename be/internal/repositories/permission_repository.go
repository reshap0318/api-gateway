package repositories

import (
	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// PermissionRepository provides database operations for Permission model.
type PermissionRepository struct {
	*GenericRepository[models.Permission]
}

// NewPermissionRepository creates a new PermissionRepository.
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		GenericRepository: NewGenericRepository(db, &models.Permission{}),
	}
}
