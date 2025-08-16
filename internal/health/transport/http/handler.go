package http

import (
	"github.com/gofiber/fiber/v2"
	health "github.com/theotruvelot/catchook/internal/health/domain"
	"github.com/theotruvelot/catchook/internal/platform/http/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
)

// Handler holds the health-specific dependencies
type Handler struct {
	healthService health.Service
}

// NewHandler creates a new health handler
func NewHandler(healthService health.Service) *Handler {
	return &Handler{
		healthService: healthService,
	}
}

func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	ctx, span := tracer.StartSpan(ctx, "health.handler")
	defer span.End()

	health, err := h.healthService.Check(ctx)
	if err != nil {
		span.RecordError(err)
		return response.InternalError(c, "health check failed")
	}
	return response.Success(c, health, "health check")
}
