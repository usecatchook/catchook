package server

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
	"strconv"
)

func (s *Server) handleGetMe(c *fiber.Ctx) error {
	currentUser := middleware.GetUser(c)
	if currentUser == nil {
		return response.Unauthorized(c, "invalid session")
	}

	userResp, err := s.container.UserService.GetByID(c.Context(), currentUser.ID)
	if err != nil {
		return response.NotFound(c, "user not found")
	}

	return response.Success(c, userResp, "user profile")
}

func (s *Server) handleGetProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	userIDStr := c.Params("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid user id", nil)
	}
	userResp, err := s.container.UserService.GetByID(ctx, userID)
	if err != nil {
		return response.NotFound(c, "user not found")
	}
	return response.Success(c, userResp, "user profile")
}

func (s *Server) handleUpdateProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	userIDStr := c.Params("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid user id", nil)
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
	ctx := c.Context()
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	users, pagination, err := s.container.UserService.List(ctx, page, limit)
	if err != nil {
		return response.InternalError(c, "failed to list users")
	}

	return response.Paginated(c, users, *pagination, "users list")
}
