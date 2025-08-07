package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
)

func (s *Server) handleHealthCheck(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	health, err := s.container.HealthService.Check(ctx)
	if err != nil {
		return response.InternalError(c, "health check failed")
	}
	return response.Success(c, health, "health check")
}
