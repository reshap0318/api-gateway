package asynq

import (
	libasynq "github.com/hibiken/asynq"

	"github.com/reshap0318/api-gateway/internal/helpers"
)

// RedisOpt builds asynq RedisClientOpt from environment variables.
// Uses the same env vars as the app Redis client for consistency.
func RedisOpt() libasynq.RedisClientOpt {
	return libasynq.RedisClientOpt{
		Addr:     helpers.GetEnv("REDIS_HOST", "localhost") + ":" + helpers.GetEnv("REDIS_PORT", "6379"),
		Username: helpers.GetEnv("REDIS_USER", ""),
		Password: helpers.GetEnv("REDIS_PASSWORD", ""),
		DB:       helpers.GetEnvInt("REDIS_DB", 0),
	}
}
