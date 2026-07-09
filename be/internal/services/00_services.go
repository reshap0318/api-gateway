package services

import (
	"time"

	"github.com/reshap0318/api-gateway/internal/pkg/email"
	"github.com/reshap0318/api-gateway/internal/database"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

// JWTConfig holds JWT configuration.
type JWTConfig struct {
	Expiration time.Duration
	RefreshExp time.Duration
}

// ServicesConfig holds all dependencies for Services.
// Add new dependencies here without changing NewServices signature.
type ServicesConfig struct {
	Repo        *repositories.Repositories
	Redis       *database.RedisCache
	Email       *email.EmailClient
	Logger      *helpers.Logger
}

// Services holds all service dependencies.
type Services struct {
	repo        *repositories.Repositories
	RedisClient *database.RedisCache
	EmailClient *email.EmailClient
	Logger      *helpers.Logger
	JWKSManager *JWKSManager
	Access      *helpers.Access
	RouteCache  RouteCacheRefresher
	cfg         *JWTConfig
}

// NewServices creates and initializes all services.
func NewServices(cfg *ServicesConfig) *Services {
	return &Services{
		repo:        cfg.Repo,
		RedisClient: cfg.Redis,
		EmailClient: cfg.Email,
		Logger:      cfg.Logger,
		cfg: &JWTConfig{
			Expiration: time.Duration(helpers.GetEnvInt("JWT_EXPIRATION", 24)) * time.Hour,
			RefreshExp: time.Duration(helpers.GetEnvInt("JWT_REFRESH_EXPIRATION", 168)) * time.Hour,
		},
	}
}
