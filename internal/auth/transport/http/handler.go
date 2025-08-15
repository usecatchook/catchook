package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/domain/auth"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/internal/platform/server"
	"github.com/theotruvelot/catchook/pkg/response"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
)

func (s *server.Server) handleLogin(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	var req auth.LoginRequest
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}
	authResp, err := s.container.AuthService.Login(ctx, req)
	if err != nil {
		return response.Unauthorized(c, "invalid credentials")
	}
	return response.Success(c, authResp, "login successful")
}

func (s *server.Server) handleRefreshSession(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	var req struct {
		SessionID string `json:"session_id" validate:"required"`
	}
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}
	sessionResp, err := s.container.AuthService.RefreshSession(ctx, req.SessionID)
	if err != nil {
		return response.Unauthorized(c, "invalid session")
	}
	return response.Success(c, sessionResp, "session refreshed")
}

func (s *server.Server) handleLogout(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	var req struct {
		SessionID string `json:"session_id" validate:"required"`
	}
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}
	if err := s.container.AuthService.Logout(ctx, req.SessionID); err != nil {
		return response.InternalError(c, "logout failed")
	}
	return response.Success(c, nil, "logout successful")
}
