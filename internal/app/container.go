package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/theotruvelot/catchook/internal/config"
	"github.com/theotruvelot/catchook/internal/domain/auth"
	"github.com/theotruvelot/catchook/internal/domain/health"
	"github.com/theotruvelot/catchook/internal/domain/setup"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/internal/repository/postgres"
	"github.com/theotruvelot/catchook/internal/service"
	"github.com/theotruvelot/catchook/pkg/cache"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/session"
	"github.com/theotruvelot/catchook/pkg/tracer"
	"github.com/theotruvelot/catchook/pkg/validator"
	postgresdb "github.com/theotruvelot/catchook/storage/postgres"
)

// Container handles dependency injection and initialization
type Container struct {
	Config    *config.Config
	AppLogger logger.Logger
	DB        *pgxpool.Pool
	Redis     *redis.Client
	Cache     cache.Cache
	Session   session.Manager
	Validator *validator.Validator

	// Services
	UserService   user.Service
	AuthService   auth.Service
	HealthService health.Service
	SetupService  setup.Service
}

// NewContainer creates and initializes all dependencies
func NewContainer(cfg *config.Config, appLogger logger.Logger) (*Container, error) {
	container := &Container{
		Config:    cfg,
		AppLogger: appLogger,
	}

	if err := container.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := container.initRedis(); err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	if err := tracer.Initialize(cfg.Tracer, appLogger); err != nil {
		appLogger.Warn(context.Background(), "Failed to initialize tracer", zap.Error(err))
	}

	container.initUtilities()
	container.initServices()

	appLogger.Info(context.Background(), "Application container initialized successfully")
	return container, nil
}

func (c *Container) initDatabase() error {
	pool, err := postgresdb.NewConnectionPool(&c.Config.Database, c.AppLogger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	c.DB = pool
	return nil
}

func (c *Container) initRedis() error {
	rdb, err := cache.NewRedisClient(&c.Config.Redis, c.AppLogger)
	if err != nil {
		return fmt.Errorf("failed to initialize redis: %w", err)
	}
	c.Redis = rdb
	return nil
}

func (c *Container) initUtilities() {
	c.Cache = cache.NewRedisCache(c.Redis)
	c.Session = session.NewManager(c.Redis, cache.TTLUserSession)
	c.Validator = validator.New()
	c.AppLogger.Info(context.Background(), "Utilities initialized")
}

func (c *Container) initServices() {
	// Repositories
	userRepo := postgres.NewUserRepository(c.DB, c.AppLogger)

	// Services
	c.UserService = service.NewUserService(userRepo, c.Cache, c.AppLogger)
	c.AuthService = service.NewAuthService(userRepo, c.Session, c.AppLogger)
	c.HealthService = service.NewHealthService(c.DB, c.Redis, userRepo, c.AppLogger, c.Config.Server.Version)
	c.SetupService = service.NewSetupService(userRepo, c.AppLogger)

	c.AppLogger.Info(context.Background(), "Services initialized")
}

// Close gracefully shuts down all connections
func (c *Container) Close() error {
	ctx := context.Background()
	c.AppLogger.Info(ctx, "Closing application connections...")

	cache.CloseRedisClient(c.Redis, c.AppLogger)
	postgresdb.ClosePool(c.DB, c.AppLogger)

	// Close tracer provider
	_ = tracer.Close(ctx)

	return nil
}
