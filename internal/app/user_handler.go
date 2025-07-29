package app

import (
	"github.com/gofiber/fiber/v2"

	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/validator"
)

func (c *Container) handleGetProfile(ctx *fiber.Ctx) error {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	c.Logger.Debug(ctx.UserContext(), "Getting user profile",
		logger.Int("user_id", userID),
	)

	foundUser, err := c.UserService.GetByID(ctx.UserContext(), userID)
	if err != nil {
		switch err {
		case user.ErrUserNotFound:
			return response.NotFound(ctx, "User not found")
		default:
			c.Logger.Error(ctx.UserContext(), "Failed to get user profile", logger.Error(err))
			return response.InternalError(ctx, "Failed to get profile")
		}
	}

	return response.Success(ctx, foundUser.ToResponse(), "Profile retrieved successfully")
}

func (c *Container) handleUpdateProfile(ctx *fiber.Ctx) error {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	var req user.UpdateRequest
	if err := c.Validator.ParseAndValidate(ctx, &req); err != nil {
		if ve, ok := err.(*validator.ValidationErrors); ok {
			return response.ValidationFailed(ctx, ve.Errors)
		}
		return response.BadRequest(ctx, "Invalid request format", nil)
	}

	c.Logger.Info(ctx.UserContext(), "Updating user profile",
		logger.Int("user_id", userID),
	)

	updatedUser, err := c.UserService.Update(ctx.UserContext(), userID, req)
	if err != nil {
		switch err {
		case user.ErrUserNotFound:
			return response.NotFound(ctx, "User not found")
		default:
			c.Logger.Error(ctx.UserContext(), "Failed to update user profile", logger.Error(err))
			return response.InternalError(ctx, "Failed to update profile")
		}
	}

	return response.Success(ctx, updatedUser.ToResponse(), "Profile updated successfully")
}

func (c *Container) handleChangePassword(ctx *fiber.Ctx) error {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	var req user.ChangePasswordRequest
	if err := c.Validator.ParseAndValidate(ctx, &req); err != nil {
		if ve, ok := err.(*validator.ValidationErrors); ok {
			return response.ValidationFailed(ctx, ve.Errors)
		}
		return response.BadRequest(ctx, "Invalid request format", nil)
	}

	c.Logger.Info(ctx.UserContext(), "Changing user password",
		logger.Int("user_id", userID),
	)

	if err := c.UserService.ChangePassword(ctx.UserContext(), userID, req); err != nil {
		switch err {
		case user.ErrUserNotFound:
			return response.NotFound(ctx, "User not found")
		case user.ErrInvalidPassword:
			return response.BadRequest(ctx, "Current password is incorrect", nil)
		default:
			c.Logger.Error(ctx.UserContext(), "Failed to change password", logger.Error(err))
			return response.InternalError(ctx, "Failed to change password")
		}
	}

	return response.Success(ctx, nil, "Password changed successfully")
}
