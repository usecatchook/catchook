package app

import (
	"github.com/gofiber/fiber/v2"

	"github.com/theotruvelot/catchook/internal/domain/auth"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/validator"
)

func (c *Container) handleRegister(ctx *fiber.Ctx) error {
	var req auth.RegisterRequest
	if err := c.Validator.ParseAndValidate(ctx, &req); err != nil {
		if ve, ok := err.(*validator.ValidationErrors); ok {
			return response.ValidationFailed(ctx, ve.Errors)
		}
		return response.BadRequest(ctx, "Invalid request format", nil)
	}

	c.Logger.Info(ctx.UserContext(), "Processing user registration",
		logger.String("email", req.Email),
	)

	authResponse, err := c.AuthService.Register(ctx.UserContext(), req)
	if err != nil {
		switch err {
		case user.ErrEmailAlreadyExists:
			return response.Conflict(ctx, "Email already exists")
		default:
			c.Logger.Error(ctx.UserContext(), "Registration failed", logger.Error(err))
			return response.InternalError(ctx, "Registration failed")
		}
	}

	return response.Created(ctx, authResponse, "User registered successfully")
}

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

func (c *Container) handleRefreshToken(ctx *fiber.Ctx) error {
	var req auth.RefreshTokenRequest
	if err := c.Validator.ParseAndValidate(ctx, &req); err != nil {
		if ve, ok := err.(*validator.ValidationErrors); ok {
			return response.ValidationFailed(ctx, ve.Errors)
		}
		return response.BadRequest(ctx, "Invalid request format", nil)
	}

	c.Logger.Debug(ctx.UserContext(), "Processing token refresh")

	tokenPair, err := c.AuthService.RefreshToken(ctx.UserContext(), req)
	if err != nil {
		switch err {
		case auth.ErrInvalidToken, auth.ErrTokenExpired:
			return response.Unauthorized(ctx, "Invalid or expired refresh token")
		default:
			c.Logger.Error(ctx.UserContext(), "Token refresh failed", logger.Error(err))
			return response.InternalError(ctx, "Token refresh failed")
		}
	}

	return response.Success(ctx, tokenPair, "Token refreshed successfully")
}
