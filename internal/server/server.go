package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/theotruvelot/catchook/internal/app"
	"github.com/theotruvelot/catchook/internal/config"
	"github.com/theotruvelot/catchook/internal/middleware"
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

func (s *Server) Shutdown() error {
	s.appLogger.Info(context.Background(), "Shutting down HTTP server...")
	return s.app.Shutdown()
}
