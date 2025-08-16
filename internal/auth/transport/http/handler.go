package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	auth "github.com/theotruvelot/catchook/internal/auth/domain"
	"github.com/theotruvelot/catchook/internal/platform/http/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
)

// Handler holds the auth-specific dependencies
type Handler struct {
	authService auth.Service
	validator   *validatorpkg.Validator
}

// NewHandler creates a new auth handler
func NewHandler(authService auth.Service, validator *validatorpkg.Validator) *Handler {
	return &Handler{
		authService: authService,
		validator:   validator,
	}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	var req auth.LoginRequest
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}
	authResp, err := h.authService.Login(ctx, req)
	if err != nil {
		return response.Unauthorized(c, "invalid credentials")
	}
	return response.Success(c, authResp, "login successful")
}

func (h *Handler) RefreshSession(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	var req struct {
		SessionID string `json:"session_id" validate:"required"`
	}
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}
	sessionResp, err := h.authService.RefreshSession(ctx, req.SessionID)
	if err != nil {
		return response.Unauthorized(c, "invalid session")
	}
	return response.Success(c, sessionResp, "session refreshed")
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	var req struct {
		SessionID string `json:"session_id" validate:"required"`
	}
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}
	if err := h.authService.Logout(ctx, req.SessionID); err != nil {
		return response.InternalError(c, "logout failed")
	}
	return response.Success(c, nil, "logout successful")
}
