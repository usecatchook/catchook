package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/platform/server"

	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
)

func (s *server.Server) handleHealthCheck(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	ctx, span := tracer.StartSpan(ctx, "health.handler")
	defer span.End()

	health, err := s.container.HealthService.Check(ctx)
	if err != nil {
		span.RecordError(err)
		return response.InternalError(c, "health check failed")
	}
	return response.Success(c, health, "health check")
}
