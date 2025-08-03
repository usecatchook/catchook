package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
)

func (c *Container) handleSetupStatus(ctx *fiber.Ctx) error {
	isFirstTime, err := c.SetupService.IsFirstTimeSetup(ctx.Context())
	if err != nil {
		c.Logger.Error(ctx.Context(), "Failed to check setup status", logger.Error(err))
		return response.InternalError(ctx, "Failed to check setup status")
	}

	return ctx.JSON(fiber.Map{
		"is_first_time_setup": isFirstTime,
		"message": func() string {
			if isFirstTime {
				return "First time setup - please create your admin account"
			}
			return "System is already configured"
		}(),
	})
}
