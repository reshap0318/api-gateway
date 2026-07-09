package repositories

import (
	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// UserRoleRepository provides database operations for UserHasRole model.
type UserRoleRepository struct {
	*GenericRepository[models.UserHasRole]
}

// NewUserRoleRepository creates a new UserRoleRepository.
func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{
		GenericRepository: NewGenericRepository(db, &models.UserHasRole{}),
	}
}

// DeleteByUserID deletes all user_has_roles for a given user ID.
func (r *UserRoleRepository) DeleteByUserID(tx *gorm.DB, userID uint) error {
	return tx.Where("user_id = ?", userID).Delete(&models.UserHasRole{}).Error
}

// FindRolesByUserID finds all roles for a given user ID.
func (r *UserRoleRepository) FindRolesByUserID(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.DB.Table("roles").
		Joins("INNER JOIN user_has_roles ON user_has_roles.role_id = roles.id").
		Where("user_has_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}
