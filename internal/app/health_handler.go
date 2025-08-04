package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
)

func (c *Container) handleHealthCheck(ctx *fiber.Ctx) error {
	status, err := c.HealthService.Check(ctx.Context())
	if err != nil {
		c.Logger.Error(ctx.Context(), "Failed to check health status", logger.Error(err))
		return response.InternalError(ctx, "Failed to check system health")
	}

	httpStatus := fiber.StatusOK
	if status.Status != "healthy" {
		httpStatus = fiber.StatusServiceUnavailable
	}

	return ctx.Status(httpStatus).JSON(status)
}
