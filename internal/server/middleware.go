package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/theotruvelot/catchook/internal/middleware"
)

// setupMiddlewares configures all middleware
func (s *Server) setupMiddlewares() {
	// Recovery middleware
	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: s.config.Logger.Development,
	}))

	// Request logging
	s.app.Use(middleware.RequestLogging(s.logger))

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
