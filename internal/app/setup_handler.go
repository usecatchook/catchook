package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/domain/setup"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/validator"
)

func (c *Container) handleSetup(ctx *fiber.Ctx) error {
	c.Logger.Info(ctx.UserContext(), "Setting up the application")

	var req setup.CreateAdminUserRequest
	if err := c.Validator.ParseAndValidate(ctx, &req); err != nil {
		if ve, ok := err.(*validator.ValidationErrors); ok {
			return response.ValidationFailed(ctx, ve.Errors)
		}
		return response.BadRequest(ctx, "Invalid request format", nil)
	}

	err := c.SetupService.CreateAdminUser(ctx.UserContext(), req)
	if err != nil {
		return response.InternalError(ctx, "Failed to create admin user")
	}

	return response.Success(ctx, nil, "Application setup completed")
}
