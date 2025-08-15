package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	otelfiber "github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/theotruvelot/catchook/internal/platform/app"

	"github.com/theotruvelot/catchook/internal/config"
	"github.com/theotruvelot/catchook/internal/platform/http/middleware"
	"github.com/theotruvelot/catchook/pkg/logger"
)

type Server struct {
	app       *fiber.App
	container *app.Container
	config    *config.Config
	appLogger logger.Logger
}

func NewServer(container *app.Container) *Server {
	server := &Server{
		container: container,
		config:    container.Config,
		appLogger: container.AppLogger,
	}

	server.app = server.createFiberApp()
	server.setupMiddlewares()
	server.setupRoutes()

	return server
}

func (s *Server) createFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		ServerHeader:          "Catchook API",
		AppName:               "Catchook API v0.0.1",
		BodyLimit:             s.config.Server.BodyLimit,
		ReadTimeout:           s.config.Server.ReadTimeout,
		WriteTimeout:          s.config.Server.WriteTimeout,
		IdleTimeout:           s.config.Server.IdleTimeout,
		DisableStartupMessage: true,
		ErrorHandler:          s.errorHandler,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})
}

func (s *Server) errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	s.appLogger.Error(middleware.GetContextWithRequestID(ctx), "HTTP error",
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

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.appLogger.Info(context.Background(), "Starting HTTP server", logger.String("address", addr))

	return s.app.Listen(addr)
}

func (s *Server) setupMiddlewares() {
	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: s.config.Logger.Development,
	}))

	s.app.Use(otelfiber.Middleware())

	s.app.Use(middleware.RequestLogging(s.appLogger))

	// Security headers
	s.app.Use(helmet.New())

	// CORS
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Compression
	s.app.Use(compress.New())

	// Rate limiting
	s.app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))
}

func (s *Server) Shutdown() error {
	s.appLogger.Info(context.Background(), "Shutting down HTTP server...")
	return s.app.Shutdown()
}
