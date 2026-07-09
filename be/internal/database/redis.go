package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig holds Redis configuration.
type RedisConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       int
}

// NewRedis creates a new Redis client.
func NewRedis(cfg RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connection established")
	return client, nil
}

// CloseRedis closes the Redis client connection.
func CloseRedis(client *redis.Client) error {
	if client != nil {
		return client.Close()
	}
	return nil
}
