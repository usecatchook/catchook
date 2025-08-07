package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/theotruvelot/catchook/internal/config"
	"github.com/theotruvelot/catchook/pkg/logger"
)

// NewRedisClient crée un nouveau client Redis
func NewRedisClient(cfg *config.RedisConfig, log logger.Logger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test de la connexion
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Info(ctx, "Redis connection established",
		logger.String("host", cfg.Host),
		logger.Int("port", cfg.Port),
		logger.Int("db", cfg.DB),
		logger.Int("pool_size", cfg.PoolSize),
	)

	return rdb, nil
}

// CloseRedisClient ferme proprement le client Redis
func CloseRedisClient(client *redis.Client, log logger.Logger) {
	if client != nil {
		if err := client.Close(); err != nil {
			log.Error(context.Background(), "Failed to close Redis connection", logger.Error(err))
		} else {
			log.Info(context.Background(), "Redis connection closed")
		}
	}
}
