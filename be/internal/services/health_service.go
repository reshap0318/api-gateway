package services

import (
	"context"
	"time"

	"github.com/reshap0318/api-gateway/internal/dtos"
)

// HealthCheckDB checks database connection.
func (s *Services) HealthCheckDB() *dtos.HealthComponent {
	start := time.Now()

	// Use SQL ping to check connection via UserRepository
	sqlDB, err := s.repo.User.DB.DB()
	if err != nil {
		return &dtos.HealthComponent{
			Status:  "disconnected",
			Message: err.Error(),
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return &dtos.HealthComponent{
			Status:  "disconnected",
			Message: err.Error(),
		}
	}

	latency := time.Since(start)
	return &dtos.HealthComponent{
		Status:  "connected",
		Latency: latency.String(),
	}
}

// HealthCheckRedis checks Redis connection (if enabled).
func (s *Services) HealthCheckRedis() *dtos.HealthComponent {
	if s.RedisClient == nil {
		return &dtos.HealthComponent{
			Status: "disabled",
		}
	}

	start := time.Now()

	if err := s.RedisClient.Ping(); err != nil {
		return &dtos.HealthComponent{
			Status:  "disconnected",
			Message: err.Error(),
		}
	}

	latency := time.Since(start)
	return &dtos.HealthComponent{
		Status:  "connected",
		Latency: latency.String(),
	}
}

// HealthGetStatus returns overall health status.
func (s *Services) HealthGetStatus(ctx context.Context) *dtos.HealthStatus {
	dbStatus := s.HealthCheckDB()
	redisStatus := s.HealthCheckRedis()

	// Determine overall status
	status := "healthy"
	if dbStatus.Status != "connected" {
		status = "unhealthy"
	} else if redisStatus.Status == "disconnected" {
		status = "degraded"
	}
	// Redis "disabled" is still considered healthy

	return &dtos.HealthStatus{
		Status:    status,
		Timestamp: time.Now(),
		Database:  dbStatus,
		Redis:     redisStatus,
	}
}
