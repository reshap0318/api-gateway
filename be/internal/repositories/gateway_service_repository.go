package repositories

import (
	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// GatewayServiceRepository provides database operations for GatewayService model.
type GatewayServiceRepository struct {
	*GenericRepository[models.GatewayService]
}

// NewGatewayServiceRepository creates a new GatewayServiceRepository.
func NewGatewayServiceRepository(db *gorm.DB) *GatewayServiceRepository {
	return &GatewayServiceRepository{
		GenericRepository: NewGenericRepository(db, &models.GatewayService{}),
	}
}

// FindAllActiveWithRoutes returns all active services with their active routes + permissions preloaded.
// Used by RouteManager to build the in-memory proxy cache.
func (r *GatewayServiceRepository) FindAllActiveWithRoutes(tx *gorm.DB) ([]models.GatewayService, error) {
	db := r.getDB(tx)
	var services []models.GatewayService
	err := db.
		Where("is_active = ?", true).
		Preload("Routes", "is_active = ?", true).
		Preload("Routes.Permissions").
		Find(&services).Error
	return services, err
}
