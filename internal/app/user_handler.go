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
	targetUserID, err := ctx.ParamsInt("id")
	if err != nil {
		return response.BadRequest(ctx, "Invalid user ID", nil)
	}

	var req user.UpdateRequest
	if err := c.Validator.ParseAndValidate(ctx, &req); err != nil {
		if ve, ok := err.(*validator.ValidationErrors); ok {
			return response.ValidationFailed(ctx, ve.Errors)
		}
		return response.BadRequest(ctx, "Invalid request format", nil)
	}

	c.Logger.Info(ctx.UserContext(), "Updating user profile",
		logger.Int("target_user_id", targetUserID),
	)

	currentUser := middleware.GetUser(ctx)
	if currentUser == nil {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	updatedUser, err := c.UserService.Update(ctx.UserContext(), targetUserID, req, currentUser)
	if err != nil {
		switch err {
		case user.ErrUserNotFound:
			return response.NotFound(ctx, "User not found")
		case user.ErrInsufficientPermissions:
			return response.Forbidden(ctx, "Insufficient permissions to update")
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

func (c *Container) handleGetMe(ctx *fiber.Ctx) error {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	foundUser, err := c.UserService.GetByID(ctx.UserContext(), userID)
	if err != nil {
		return response.InternalError(ctx, "Failed to get user")
	}

	return response.Success(ctx, foundUser.ToResponse(), "User retrieved successfully")
}

func (c *Container) handleListUsers(ctx *fiber.Ctx) error {
	// Récupération des paramètres de pagination depuis les query parameters
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	// Validation des paramètres
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	c.Logger.Debug(ctx.UserContext(), "Listing users with pagination",
		logger.Int("page", page),
		logger.Int("limit", limit),
	)

	users, pagination, err := c.UserService.List(ctx.UserContext(), page, limit)
	if err != nil {
		c.Logger.Error(ctx.UserContext(), "Failed to list users", logger.Error(err))
		return response.InternalError(ctx, "Failed to list users")
	}

	return response.Paginated(ctx, users, *pagination, "Users retrieved successfully")
}
