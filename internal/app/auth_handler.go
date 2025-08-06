package app

import (
	"github.com/gofiber/fiber/v2"

	"github.com/theotruvelot/catchook/internal/domain/auth"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/validator"
)

func (c *Container) handleLogin(ctx *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.Validator.ParseAndValidate(ctx, &req); err != nil {
		if ve, ok := err.(*validator.ValidationErrors); ok {
			return response.ValidationFailed(ctx, ve.Errors)
		}
		return response.BadRequest(ctx, "Invalid request format", nil)
	}

	c.Logger.Info(ctx.UserContext(), "Processing user login",
		logger.String("email", req.Email),
	)

	authResponse, err := c.AuthService.Login(ctx.UserContext(), req)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			return response.Unauthorized(ctx, "Invalid credentials")
		case user.ErrUserInactive:
			return response.Forbidden(ctx, "Account is inactive")
		default:
			c.Logger.Error(ctx.UserContext(), "Login failed", logger.Error(err))
			return response.InternalError(ctx, "Login failed")
		}
	}

	return response.Success(ctx, authResponse, "Login successful")
}

func (c *Container) handleRefreshSession(ctx *fiber.Ctx) error {
	sessionID := ctx.Get("Authorization")
	if sessionID == "" {
		return response.Unauthorized(ctx, "Missing session ID")
	}

	c.Logger.Debug(ctx.UserContext(), "Processing session refresh")

	sessionResponse, err := c.AuthService.RefreshSession(ctx.UserContext(), sessionID)
	if err != nil {
		switch err {
		case auth.ErrInvalidToken:
			return response.Unauthorized(ctx, "Invalid or expired session")
		default:
			c.Logger.Error(ctx.UserContext(), "Session refresh failed", logger.Error(err))
			return response.InternalError(ctx, "Session refresh failed")
		}
	}

	return response.Success(ctx, sessionResponse, "Session refreshed successfully")
}

func (c *Container) handleLogout(ctx *fiber.Ctx) error {
	sessionID := ctx.Get("Authorization")
	if sessionID == "" {
		return response.Unauthorized(ctx, "Missing session ID")
	}

	c.Logger.Debug(ctx.UserContext(), "Processing logout")

	err := c.AuthService.Logout(ctx.UserContext(), sessionID)
	if err != nil {
		c.Logger.Error(ctx.UserContext(), "Logout failed", logger.Error(err))
		return response.InternalError(ctx, "Logout failed")
	}

	return response.Success(ctx, nil, "Logout successful")
}
