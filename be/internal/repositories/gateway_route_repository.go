package repositories

import (
	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// GatewayRouteRepository provides database operations for GatewayRoute model.
type GatewayRouteRepository struct {
	*GenericRepository[models.GatewayRoute]
}

// NewGatewayRouteRepository creates a new GatewayRouteRepository.
func NewGatewayRouteRepository(db *gorm.DB) *GatewayRouteRepository {
	return &GatewayRouteRepository{
		GenericRepository: NewGenericRepository(db, &models.GatewayRoute{}),
	}
}
