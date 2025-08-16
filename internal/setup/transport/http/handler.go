package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/platform/http/middleware"
	setup "github.com/theotruvelot/catchook/internal/setup/domain"
	"github.com/theotruvelot/catchook/pkg/response"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
)

// Handler holds the setup-specific dependencies
type Handler struct {
	setupService setup.Service
	validator    *validatorpkg.Validator
}

// NewHandler creates a new setup handler
func NewHandler(setupService setup.Service, validator *validatorpkg.Validator) *Handler {
	return &Handler{
		setupService: setupService,
		validator:    validator,
	}
}

func (h *Handler) Setup(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	var req setup.CreateAdminUserRequest
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}
	if err := h.setupService.CreateAdminUser(ctx, req); err != nil {
		if errors.Is(err, setup.ErrAdminUserAlreadyExists) {
			return response.Conflict(c, "admin user already exists")
		}
		return response.InternalError(c, "setup failed")
	}
	return response.Success(c, nil, "setup completed")
}
