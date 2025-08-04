package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/theotruvelot/catchook/internal/domain/health"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/pkg/logger"
)

type healthService struct {
	db       *pgxpool.Pool
	redis    *redis.Client
	userRepo user.Repository
	logger   logger.Logger
	version  string
}

func NewHealthService(
	db *pgxpool.Pool,
	redis *redis.Client,
	userRepo user.Repository,
	logger logger.Logger,
	version string,
) health.Service {
	return &healthService{
		db:       db,
		redis:    redis,
		userRepo: userRepo,
		logger:   logger,
		version:  version,
	}
}

func (s *healthService) Check(ctx context.Context) (*health.StatusResponse, error) {
	services := make(map[string]string)

	// Check database
	services["database"] = "ok"
	if err := s.db.Ping(ctx); err != nil {
		services["database"] = "error"
		s.logger.Error(ctx, "Database health check failed", logger.Error(err))
	}

	// Check Redis
	services["redis"] = "ok"
	if err := s.redis.Ping(ctx).Err(); err != nil {
		services["redis"] = "error"
		s.logger.Error(ctx, "Redis health check failed", logger.Error(err))
	}

	// Check if first time setup
	isFirstTime := false
	count, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to count users during setup check", logger.Error(err))
	} else {
		isFirstTime = count == 0
	}

	message := "System is operational"
	if isFirstTime {
		message = "System requires initial setup"
	}

	status := "healthy"
	if services["database"] == "error" || services["redis"] == "error" {
		status = "unhealthy"
	}

	return &health.StatusResponse{
		Status:         status,
		Version:        s.version,
		FirstTimeSetup: isFirstTime,
		Message:        message,
		Services:       services,
	}, nil
}
