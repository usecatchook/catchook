package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

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
	"github.com/theotruvelot/catchook/pkg/validator"
	postgresdb "github.com/theotruvelot/catchook/storage/postgres"
)

// Container handles dependency injection and initialization
type Container struct {
	Config    *config.Config
	Logger    logger.Logger
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
func NewContainer(cfg *config.Config, logger logger.Logger) (*Container, error) {
	container := &Container{
		Config: cfg,
		Logger: logger,
	}

	if err := container.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := container.initRedis(); err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	container.initUtilities()
	container.initServices()

	logger.Info(context.Background(), "Application container initialized successfully")
	return container, nil
}

func (c *Container) initDatabase() error {
	pool, err := postgresdb.NewConnectionPool(&c.Config.Database, c.Logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	c.DB = pool
	return nil
}

func (c *Container) initRedis() error {
	rdb, err := cache.NewRedisClient(&c.Config.Redis, c.Logger)
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
	c.Logger.Info(context.Background(), "Utilities initialized")
}

func (c *Container) initServices() {
	// Repositories
	userRepo := postgres.NewUserRepository(c.DB, c.Logger)

	// Services
	c.UserService = service.NewUserService(userRepo, c.Cache, c.Logger)
	c.AuthService = service.NewAuthService(userRepo, c.Session, c.Logger)
	c.HealthService = service.NewHealthService(c.DB, c.Redis, userRepo, c.Logger, c.Config.Server.Version)
	c.SetupService = service.NewSetupService(userRepo, c.Logger)

	c.Logger.Info(context.Background(), "Services initialized")
}

// Close gracefully shuts down all connections
func (c *Container) Close() error {
	ctx := context.Background()
	c.Logger.Info(ctx, "Closing application connections...")

	cache.CloseRedisClient(c.Redis, c.Logger)
	postgresdb.ClosePool(c.DB, c.Logger)

	return nil
}
