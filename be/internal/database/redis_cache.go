package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache provides helper functions for Redis caching operations.
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new RedisCache instance.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

// Get retrieves a value from Redis by key.
func (r *RedisCache) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// Set sets a value in Redis with optional expiration.
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Delete deletes one or more keys from Redis.
func (r *RedisCache) Delete(keys ...string) error {
	return r.client.Del(r.ctx, keys...).Err()
}

// Exists checks if a key exists in Redis.
func (r *RedisCache) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Expire sets expiration time for a key.
func (r *RedisCache) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// TTL gets the time to live for a key.
func (r *RedisCache) TTL(key string) (time.Duration, error) {
	return r.client.TTL(r.ctx, key).Result()
}

// GetJSON retrieves a JSON value from Redis and unmarshals it into the provided struct.
func (r *RedisCache) GetJSON(key string, dest interface{}) error {
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// SetJSON marshals a value to JSON and stores it in Redis.
func (r *RedisCache) SetJSON(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return r.client.Set(r.ctx, key, data, expiration).Err()
}

// Incr increments a key's value by 1.
func (r *RedisCache) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// IncrBy increments a key's value by the given amount.
func (r *RedisCache) IncrBy(key string, amount int64) (int64, error) {
	return r.client.IncrBy(r.ctx, key, amount).Result()
}

// Decr decrements a key's value by 1.
func (r *RedisCache) Decr(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

// DecrBy decrements a key's value by the given amount.
func (r *RedisCache) DecrBy(key string, amount int64) (int64, error) {
	return r.client.DecrBy(r.ctx, key, amount).Result()
}

// HSet sets a hash field to a value.
func (r *RedisCache) HSet(key, field string, value interface{}) error {
	return r.client.HSet(r.ctx, key, field, value).Err()
}

// HGet gets a hash field value.
func (r *RedisCache) HGet(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

// HGetAll gets all hash fields and values.
func (r *RedisCache) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

// HDel deletes one or more hash fields.
func (r *RedisCache) HDel(key string, fields ...string) error {
	return r.client.HDel(r.ctx, key, fields...).Err()
}

// HExists checks if a hash field exists.
func (r *RedisCache) HExists(key, field string) (bool, error) {
	return r.client.HExists(r.ctx, key, field).Result()
}

// LPush pushes a value to the left of a list.
func (r *RedisCache) LPush(key string, values ...interface{}) error {
	return r.client.LPush(r.ctx, key, values...).Err()
}

// RPush pushes a value to the right of a list.
func (r *RedisCache) RPush(key string, values ...interface{}) error {
	return r.client.RPush(r.ctx, key, values...).Err()
}

// LRange gets a range of elements from a list.
func (r *RedisCache) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(r.ctx, key, start, stop).Result()
}

// LLen gets the length of a list.
func (r *RedisCache) LLen(key string) (int64, error) {
	return r.client.LLen(r.ctx, key).Result()
}

// SAdd adds members to a set.
func (r *RedisCache) SAdd(key string, members ...interface{}) error {
	return r.client.SAdd(r.ctx, key, members...).Err()
}

// SMembers gets all members of a set.
func (r *RedisCache) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}

// SIsMember checks if a value is a member of a set.
func (r *RedisCache) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(r.ctx, key, member).Result()
}

// SRem removes members from a set.
func (r *RedisCache) SRem(key string, members ...interface{}) error {
	return r.client.SRem(r.ctx, key, members...).Err()
}

// ZAdd adds a member to a sorted set with a score.
func (r *RedisCache) ZAdd(key string, score float64, member interface{}) error {
	return r.client.ZAdd(r.ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// ZRange gets a range of members from a sorted set.
func (r *RedisCache) ZRange(key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(r.ctx, key, start, stop).Result()
}

// ZRem removes a member from a sorted set.
func (r *RedisCache) ZRem(key string, members ...interface{}) error {
	return r.client.ZRem(r.ctx, key, members...).Err()
}

// FlushDB flushes the current database.
func (r *RedisCache) FlushDB() error {
	return r.client.FlushDB(r.ctx).Err()
}

// Ping checks Redis connection.
func (r *RedisCache) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// IsCacheAvailable checks if cache is enabled and available
func (r *RedisCache) IsCacheAvailable() bool {
	return r.client != nil
}

// Publish publishes a message to a Redis Pub/Sub channel.
func (r *RedisCache) Publish(channel string, message interface{}) error {
	return r.client.Publish(r.ctx, channel, message).Err()
}

// Subscribe subscribes to a Redis Pub/Sub channel. The caller is responsible for
// reading from the returned PubSub's Channel() and calling Close() when done.
func (r *RedisCache) Subscribe(channel string) *redis.PubSub {
	return r.client.Subscribe(r.ctx, channel)
}