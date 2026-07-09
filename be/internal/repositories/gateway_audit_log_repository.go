package repositories

import (
	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// GatewayAuditLogRepository provides database operations for GatewayAuditLog model.
type GatewayAuditLogRepository struct {
	*GenericRepository[models.GatewayAuditLog]
}

// NewGatewayAuditLogRepository creates a new GatewayAuditLogRepository.
func NewGatewayAuditLogRepository(db *gorm.DB) *GatewayAuditLogRepository {
	return &GatewayAuditLogRepository{
		GenericRepository: NewGenericRepository(db, &models.GatewayAuditLog{}),
	}
}
