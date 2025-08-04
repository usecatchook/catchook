package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/theotruvelot/catchook/internal/config"
	"github.com/theotruvelot/catchook/internal/domain/auth"
	"github.com/theotruvelot/catchook/internal/domain/health"
	"github.com/theotruvelot/catchook/internal/domain/setup"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/internal/repository/postgres"
	"github.com/theotruvelot/catchook/internal/service"
	"github.com/theotruvelot/catchook/pkg/cache"
	"github.com/theotruvelot/catchook/pkg/jwt"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/validator"
)

type Container struct {
	Config    *config.Config
	Logger    logger.Logger
	DB        *pgxpool.Pool
	Redis     *redis.Client
	Cache     cache.Cache
	JWT       jwt.Manager
	Validator *validator.Validator

	UserService   user.Service
	AuthService   auth.Service
	HealthService health.Service
	SetupService  setup.Service
}

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
	config, err := pgxpool.ParseConfig(c.Config.Database.DatabaseURL())
	if err != nil {
		return fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConns = int32(c.Config.Database.MaxOpenConns)
	config.MinConns = int32(c.Config.Database.MaxIdleConns)
	config.MaxConnLifetime = c.Config.Database.ConnMaxLifetime
	config.MaxConnIdleTime = c.Config.Database.ConnMaxIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	c.DB = pool
	c.Logger.Info(context.Background(), "Database connection established",
		logger.String("host", c.Config.Database.Host),
		logger.Int("port", c.Config.Database.Port),
		logger.String("database", c.Config.Database.Name),
	)

	return nil
}

func (c *Container) initUtilities() {
	c.Cache = cache.NewRedisCache(c.Redis)
	c.JWT = jwt.NewManager(c.Config.JWT)
	c.Validator = validator.New()
	c.Logger.Info(context.Background(), "Utilities initialized")
}

func (c *Container) initServices() {
	// Repositories
	userRepo := postgres.NewUserRepository(c.DB, c.Logger)

	// Services
	c.UserService = service.NewUserService(userRepo, c.Cache, c.Logger)
	c.AuthService = service.NewAuthService(userRepo, c.UserService, c.JWT, c.Logger)
	c.HealthService = service.NewHealthService(c.DB, c.Redis, userRepo, c.Logger, c.Config.Server.Version)
	c.SetupService = service.NewSetupService(userRepo, c.Logger)

	c.Logger.Info(context.Background(), "Services initialized")
}

func (c *Container) initRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Config.Redis.RedisAddr(),
		Password:     c.Config.Redis.Password,
		DB:           c.Config.Redis.DB,
		PoolSize:     c.Config.Redis.PoolSize,
		MinIdleConns: c.Config.Redis.MinIdleConns,
		DialTimeout:  c.Config.Redis.DialTimeout,
		ReadTimeout:  c.Config.Redis.ReadTimeout,
		WriteTimeout: c.Config.Redis.WriteTimeout,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	c.Redis = rdb
	c.Logger.Info(ctx, "Redis connection established",
		logger.String("host", c.Config.Redis.Host),
		logger.Int("port", c.Config.Redis.Port),
		logger.Int("db", c.Config.Redis.DB),
	)

	return nil
}

func (c *Container) NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ServerHeader:          "Catchook API",
		AppName:               "Catchook API v0.0.1",
		BodyLimit:             c.Config.Server.BodyLimit,
		ReadTimeout:           c.Config.Server.ReadTimeout,
		WriteTimeout:          c.Config.Server.WriteTimeout,
		IdleTimeout:           c.Config.Server.IdleTimeout,
		DisableStartupMessage: true,
		ErrorHandler:          c.errorHandler,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	c.setupMiddlewares(app)

	return app
}

func (c *Container) setupMiddlewares(app *fiber.App) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: c.Config.Logger.Development,
	}))

	app.Use(middleware.RequestLogging(c.Logger))

	app.Use(helmet.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Use(compress.New())

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))
}

func (c *Container) SetupRoutes(app *fiber.App) {
	// Health check
	app.Get("/health", c.handleHealthCheck)

	// API routes
	api := app.Group("/api/v1")

	c.setupAuthRoutes(api)

	c.setupUserRoutes(api)

	c.setupSetupRoutes(api)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return response.NotFound(c, "Route not found")
	})
}

func (c *Container) setupAuthRoutes(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/login", c.handleLogin)
	auth.Post("/refresh", c.handleRefreshToken)
}

func (c *Container) setupUserRoutes(api fiber.Router) {
	users := api.Group("/users")

	users.Use(middleware.JWTAuth(c.JWT))

	users.Get("/profile/:id", c.handleGetProfile)
	users.Put("/profile/:id", c.handleUpdateProfile)
	users.Post("/change-password", c.handleChangePassword)
}

func (c *Container) setupSetupRoutes(api fiber.Router) {
	setup := api.Group("/setup")

	setup.Post("/", c.handleSetup)
}

func (c *Container) errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	c.Logger.Error(ctx.Context(), "HTTP error",
		logger.String("method", ctx.Method()),
		logger.String("path", ctx.Path()),
		logger.String("ip", ctx.IP()),
		logger.Int("status", code),
		logger.Error(err),
	)

	return ctx.Status(code).JSON(fiber.Map{
		"error": message,
		"code":  fmt.Sprintf("HTTP_%d", code),
	})
}

func (c *Container) Close() error {
	ctx := context.Background()
	c.Logger.Info(ctx, "Closing application connections...")

	if c.Redis != nil {
		if err := c.Redis.Close(); err != nil {
			c.Logger.Error(ctx, "Failed to close Redis connection", logger.Error(err))
		}
	}

	if c.DB != nil {
		c.DB.Close()
	}

	return nil
}
