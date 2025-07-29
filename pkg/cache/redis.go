package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetJSON(ctx context.Context, key string, dest interface{}) error
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	FlushAll(ctx context.Context) error
	Pipeline() Pipeline
}

type Pipeline interface {
	Set(key string, value interface{}, ttl time.Duration) Pipeline
	Get(key string) Pipeline
	Delete(keys ...string) Pipeline
	Exec(ctx context.Context) error
}

type redisCache struct {
	client *redis.Client
}

type redisPipeline struct {
	pipe redis.Pipeliner
	ops  []func() error
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
	if err == redis.Nil {
		return "", fmt.Errorf("key not found")
	}
	return result, err
}

func (r *redisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
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

func (r *redisCache) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *redisCache) Decrement(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

func (r *redisCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, ttl).Result()
}

func (r *redisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// FlushAll supprime toutes les clés
func (r *redisCache) FlushAll(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

func (r *redisCache) Pipeline() Pipeline {
	return &redisPipeline{
		pipe: r.client.Pipeline(),
		ops:  make([]func() error, 0),
	}
}

func (rp *redisPipeline) Set(key string, value interface{}, ttl time.Duration) Pipeline {
	rp.ops = append(rp.ops, func() error {
		return rp.pipe.Set(context.Background(), key, value, ttl).Err()
	})
	return rp
}

func (rp *redisPipeline) Get(key string) Pipeline {
	rp.ops = append(rp.ops, func() error {
		return rp.pipe.Get(context.Background(), key).Err()
	})
	return rp
}

func (rp *redisPipeline) Delete(keys ...string) Pipeline {
	rp.ops = append(rp.ops, func() error {
		return rp.pipe.Del(context.Background(), keys...).Err()
	})
	return rp
}

func (rp *redisPipeline) Exec(ctx context.Context) error {
	_, err := rp.pipe.Exec(ctx)
	return err
}

func GetOrSet[T any](ctx context.Context, cache Cache, key string, ttl time.Duration, fetchFunc func() (T, error)) (T, error) {
	var result T

	err := cache.GetJSON(ctx, key, &result)
	if err == nil {
		return result, nil
	}

	result, err = fetchFunc()
	if err != nil {
		return result, err
	}

	cache.SetJSON(ctx, key, result, ttl)

	return result, nil
}

type prefixedCache struct {
	cache  Cache
	prefix string
}

func NewPrefixedCache(cache Cache, prefix string) Cache {
	return &prefixedCache{
		cache:  cache,
		prefix: prefix + ":",
	}
}

func (p *prefixedCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return p.cache.Set(ctx, p.prefix+key, value, ttl)
}

func (p *prefixedCache) Get(ctx context.Context, key string) (string, error) {
	return p.cache.Get(ctx, p.prefix+key)
}

func (p *prefixedCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	return p.cache.GetJSON(ctx, p.prefix+key, dest)
}

func (p *prefixedCache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return p.cache.SetJSON(ctx, p.prefix+key, value, ttl)
}

func (p *prefixedCache) Delete(ctx context.Context, keys ...string) error {
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = p.prefix + key
	}
	return p.cache.Delete(ctx, prefixedKeys...)
}

func (p *prefixedCache) Exists(ctx context.Context, key string) (bool, error) {
	return p.cache.Exists(ctx, p.prefix+key)
}

func (p *prefixedCache) Increment(ctx context.Context, key string) (int64, error) {
	return p.cache.Increment(ctx, p.prefix+key)
}

func (p *prefixedCache) Decrement(ctx context.Context, key string) (int64, error) {
	return p.cache.Decrement(ctx, p.prefix+key)
}

func (p *prefixedCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	return p.cache.SetNX(ctx, p.prefix+key, value, ttl)
}

func (p *prefixedCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return p.cache.TTL(ctx, p.prefix+key)
}

func (p *prefixedCache) FlushAll(ctx context.Context) error {
	return p.cache.FlushAll(ctx)
}

func (p *prefixedCache) Pipeline() Pipeline {
	return p.cache.Pipeline()
}

const (
	KeyUserSession     = "user:session:%d"     // user:session:123
	KeyUserProfile     = "user:profile:%d"     // user:profile:123
	KeyWebhookDelivery = "webhook:delivery:%s" // webhook:delivery:uuid
	KeyRateLimit       = "rate:limit:%s:%s"    // rate:limit:ip:endpoint
	KeyEmailVerify     = "email:verify:%s"     // email:verify:token
	KeyPasswordReset   = "password:reset:%s"   // password:reset:token
)

func BuildKey(template string, args ...interface{}) string {
	return fmt.Sprintf(template, args...)
}

const (
	TTLUserSession     = 24 * time.Hour
	TTLUserProfile     = 1 * time.Hour
	TTLWebhookDelivery = 7 * 24 * time.Hour
	TTLRateLimit       = 1 * time.Minute
	TTLEmailVerify     = 15 * time.Minute
	TTLPasswordReset   = 15 * time.Minute
	TTLShortTerm       = 5 * time.Minute
	TTLMediumTerm      = 1 * time.Hour
	TTLLongTerm        = 24 * time.Hour
)
