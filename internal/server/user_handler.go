package server

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
	"github.com/theotruvelot/catchook/storage/postgres/generated"
)

func (s *Server) handleGetMe(c *fiber.Ctx) error {
	currentUser := middleware.GetUser(c)
	if currentUser == nil {
		return response.Unauthorized(c, "invalid session")
	}

	ctx := middleware.GetContextWithRequestID(c)
	userResp, err := s.container.UserService.GetByID(ctx, currentUser.ID)
	if err != nil {
		return response.NotFound(c, "user not found")
	}

	return response.Success(c, userResp, "user profile")
}

func (s *Server) handleGetProfile(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	userID := c.Params("id")
	if userID == "" {
		return response.BadRequest(c, "user_id is required", nil)
	}
	userResp, err := s.container.UserService.GetByID(ctx, userID)
	if err != nil {
		return response.NotFound(c, "user not found")
	}
	return response.Success(c, userResp, "user profile")
}

func (s *Server) handleUpdateProfile(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	userID := c.Params("id")
	if userID == "" {
		return response.BadRequest(c, "user_id is required", nil)
	}

	currentUser := middleware.GetUser(c)
	if currentUser == nil {
		return response.Unauthorized(c, "invalid session")
	}
	if currentUser.ID != userID && currentUser.Role != "admin" {
		return response.Forbidden(c, "can only update own profile")
	}
	var req user.UpdateRequest
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}
	userResp, err := s.container.UserService.Update(ctx, userID, req, currentUser)
	if err != nil {
		return response.InternalError(c, "update failed")
	}
	return response.Success(c, userResp, "profile updated")
}

func (s *Server) handleListUsers(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	users, pagination, err := s.container.UserService.List(ctx, page, limit)
	if err != nil {
		return response.InternalError(c, "failed to list users")
	}

	return response.Paginated(c, users, *pagination, "users list")
}

func (s *Server) handleCreateUser(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	currentUser := middleware.GetUser(c)
	if currentUser == nil || currentUser.Role != generated.UserRoleAdmin {
		return response.Forbidden(c, "only admins can create users")
	}

	var req user.CreateRequest
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	newUser, err := s.container.UserService.Create(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailAlreadyExists):
			return response.Conflict(c, "email already exists")
		default:
			return response.InternalError(c, "failed to create user")
		}
	}

	return response.Success(c, newUser, "user created")
}
