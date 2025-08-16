package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/platform/http/middleware"
	user "github.com/theotruvelot/catchook/internal/user/domain"
	"github.com/theotruvelot/catchook/pkg/response"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
)

// Handler holds the user-specific dependencies
type Handler struct {
	userService user.Service
	validator   *validatorpkg.Validator
}

// NewHandler creates a new user handler
func NewHandler(userService user.Service, validator *validatorpkg.Validator) *Handler {
	return &Handler{
		userService: userService,
		validator:   validator,
	}
}

func (h *Handler) GetMe(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	currentUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return response.Unauthorized(c, "invalid session")
	}

	userResp, err := h.userService.GetByID(ctx, currentUser.ID)
	if err != nil {
		return response.NotFound(c, "user not found")
	}

	return response.Success(c, userResp.ToResponse(), "user profile")
}

func (h *Handler) GetProfile(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	userID := c.Params("id")
	if userID == "" {
		return response.BadRequest(c, "user_id is required", nil)
	}
	userResp, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return response.NotFound(c, "user not found")
		}
		return response.InternalError(c, "failed to get user profile")
	}
	return response.Success(c, userResp.ToResponse(), "user profile")
}

func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	userID := c.Params("id")
	if userID == "" {
		return response.BadRequest(c, "user_id is required", nil)
	}

	currentUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return response.Unauthorized(c, "invalid session")
	}

	var req user.UpdateRequest
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	// Convert middleware User to domain CurrentUser
	domainUser := &user.CurrentUser{
		ID:   currentUser.ID,
		Role: string(currentUser.Role),
	}

	userResp, err := h.userService.Update(ctx, userID, req, domainUser)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return response.NotFound(c, "user not found")
		}
		if errors.Is(err, user.ErrInsufficientPermissions) {
			return response.Forbidden(c, "insufficient permissions")
		}
		return response.InternalError(c, "update failed")
	}
	return response.Success(c, userResp.ToResponse(), "profile updated")
}

func (h *Handler) ListUsers(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	users, pagination, err := h.userService.List(ctx, page, limit)
	if err != nil {
		return response.InternalError(c, "failed to list users")
	}

	// Convert to response format
	userResponses := make([]*user.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = u.ToResponse()
	}

	listResp := &user.ListUsersResponse{
		Users:      userResponses,
		Pagination: pagination,
	}

	return response.Success(c, listResp, "users listed")
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	var req user.CreateRequest
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	newUser, err := h.userService.Create(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailAlreadyExists):
			return response.Conflict(c, "email already exists")
		default:
			return response.InternalError(c, "failed to create user")
		}
	}

	return response.Success(c, newUser.ToResponse(), "user created")
}
