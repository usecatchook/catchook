package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/theotruvelot/catchook/internal/app"
	"github.com/theotruvelot/catchook/internal/config"
	"github.com/theotruvelot/catchook/pkg/logger"
)

// Server handles HTTP server configuration and setup
type Server struct {
	app       *fiber.App
	container *app.Container
	config    *config.Config
	logger    logger.Logger
}

// NewServer creates a new HTTP server with all configurations
func NewServer(container *app.Container) *Server {
	server := &Server{
		container: container,
		config:    container.Config,
		logger:    container.Logger,
	}

	server.app = server.createFiberApp()
	server.setupMiddlewares()
	server.setupRoutes()

	return server
}

// createFiberApp creates and configures the Fiber app
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

// errorHandler handles all HTTP errors
func (s *Server) errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	s.logger.Error(ctx.Context(), "HTTP error",
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

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.logger.Info(context.Background(), "Starting HTTP server", logger.String("address", addr))

	return s.app.Listen(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	s.logger.Info(context.Background(), "Shutting down HTTP server...")
	return s.app.Shutdown()
}

// App returns the Fiber app instance (for testing)
func (s *Server) App() *fiber.App {
	return s.app
}
