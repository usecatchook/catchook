package server

import (
	"time"

	otelfiber "github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/theotruvelot/catchook/internal/middleware"
)

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
