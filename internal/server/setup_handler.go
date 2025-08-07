package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/domain/setup"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
)

func (s *Server) handleSetup(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	var req setup.CreateAdminUserRequest
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}
	if err := s.container.SetupService.CreateAdminUser(ctx, req); err != nil {
		return response.InternalError(c, "setup failed")
	}
	return response.Success(c, nil, "setup completed")
}
