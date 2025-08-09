package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/theotruvelot/catchook/internal/domain/health"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span := tracer.StartSpan(ctx, "health.service.check")
	defer span.End()

	services := make(map[string]string)

	services["database"] = "ok"
	if err := tracer.WithSpan(ctx, "health.db.ping", func(inner context.Context) error { return s.db.Ping(inner) }); err != nil {
		services["database"] = "error"
		s.logger.Error(ctx, "Database health check failed", logger.Error(err))
		span.RecordError(err)
	}

	services["redis"] = "ok"
	if err := tracer.WithSpan(ctx, "health.redis.ping", func(inner context.Context) error { return s.redis.Ping(inner).Err() }); err != nil {
		services["redis"] = "error"
		s.logger.Error(ctx, "Redis health check failed", logger.Error(err))
		span.RecordError(err)
	}

	isFirstTime := false
	var count int64
	if err := tracer.WithSpan(ctx, "health.service.count_users", func(inner context.Context) error {
		var err error
		count, err = s.userRepo.CountUsers(inner)
		return err
	}); err != nil {
		s.logger.Error(ctx, "Failed to count users during setup check", logger.Error(err))
		span.RecordError(err)
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

	resp := &health.StatusResponse{
		Status:         status,
		Version:        s.version,
		FirstTimeSetup: isFirstTime,
		Message:        message,
		Services:       services,
	}

	span.SetAttributes(
		attribute.String("health.status", status),
		attribute.Bool("health.first_time_setup", isFirstTime),
	)

	return resp, nil
}
