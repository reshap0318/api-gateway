package repositories

import (
	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// RoleRepository provides database operations for Role model.
type RoleRepository struct {
	*GenericRepository[models.Role]
}

// NewRoleRepository creates a new RoleRepository.
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		GenericRepository: NewGenericRepository(db, &models.Role{}),
	}
}
