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

// GetCacheVersion returns the current Service/Route CUD version (gateway_cache_meta,
// singleton row id=1) — a cheap check so RouteManager's periodic refresh can skip
// rebuilding when nothing changed.
func (r *GatewayServiceRepository) GetCacheVersion(tx *gorm.DB) (uint64, error) {
	db := r.getDB(tx)
	var meta models.GatewayCacheMeta
	if err := db.First(&meta, 1).Error; err != nil {
		return 0, err
	}
	return meta.Version, nil
}

// TouchCacheVersion atomically increments gateway_cache_meta's version — a plain DB-level
// `version + 1`, not read-then-write, so concurrent CUD calls never lose an increment to a race.
func (r *GatewayServiceRepository) TouchCacheVersion(tx *gorm.DB) error {
	db := r.getDB(tx)
	return db.Model(&models.GatewayCacheMeta{}).Where("id = ?", 1).
		Update("version", gorm.Expr("version + 1")).Error
}
