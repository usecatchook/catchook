package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theotruvelot/catchook/internal/config"
	"github.com/theotruvelot/catchook/pkg/logger"
)

func NewConnectionPool(cfg *config.DatabaseConfig, log logger.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configuration du pool de connexions
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configuration des paramètres du pool
	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = cfg.ConnMaxIdleTime

	// Configuration des timeouts
	poolConfig.ConnConfig.ConnectTimeout = 10 * time.Second
	poolConfig.ConnConfig.RuntimeParams["application_name"] = "catchook-api"

	// Création du pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test de la connexion
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info(ctx, "PostgreSQL connection pool established",
		logger.Int("max_connections", cfg.MaxOpenConns),
		logger.Int("max_idle_connections", cfg.MaxIdleConns),
		logger.Duration("max_lifetime", cfg.ConnMaxLifetime),
	)

	return pool, nil
}

// ClosePool ferme proprement le pool de connexions
func ClosePool(pool *pgxpool.Pool, log logger.Logger) {
	if pool != nil {
		pool.Close()
		log.Info(context.Background(), "PostgreSQL connection pool closed")
	}
}
