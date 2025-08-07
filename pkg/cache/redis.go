package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache interface simplifiée avec seulement les opérations CRUD de base
type Cache interface {
	// Create
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Read
	Get(ctx context.Context, key string) (string, error)
	GetJSON(ctx context.Context, key string, dest interface{}) error
	Exists(ctx context.Context, key string) (bool, error)

	// Update
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)

	// Delete
	Delete(ctx context.Context, keys ...string) error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{
		client: client,
	}
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("key not found")
	}
	return result, err
}

func (r *redisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	result, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("key not found")
	}
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(result), dest)
}

func (r *redisCache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return r.client.Set(ctx, key, jsonData, ttl).Err()
}

func (r *redisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

func (r *redisCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, ttl).Result()
}

const (
	KeyUserSession = "user:session:%s"
	KeyUserProfile = "user:profile:%s"
)

func BuildKey(template string, args ...interface{}) string {
	return fmt.Sprintf(template, args...)
}

const (
	TTLUserSession = 24 * time.Hour
	TTLUserProfile = 1 * time.Hour
)
