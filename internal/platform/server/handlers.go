package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	auth "github.com/theotruvelot/catchook/internal/auth/domain"
	"github.com/theotruvelot/catchook/internal/platform/http/middleware"
	setup "github.com/theotruvelot/catchook/internal/setup/domain"
	source "github.com/theotruvelot/catchook/internal/source/domain"
	user "github.com/theotruvelot/catchook/internal/user/domain"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
)

// Health handler
func (s *Server) handleHealthCheck(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	ctx, span := tracer.StartSpan(ctx, "health.handler")
	defer span.End()

	health, err := s.container.HealthService.Check(ctx)
	if err != nil {
		span.RecordError(err)
		return response.InternalError(c, "health check failed")
	}
	return response.Success(c, health, "health check")
}

// Auth handlers
func (s *Server) handleLogin(c *fiber.Ctx) error {
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

func (s *Server) handleRefreshSession(c *fiber.Ctx) error {
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

func (s *Server) handleLogout(c *fiber.Ctx) error {
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

// Setup handlers
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
		if errors.Is(err, setup.ErrAdminUserAlreadyExists) {
			return response.Conflict(c, "admin user already exists")
		}
		return response.InternalError(c, "setup failed")
	}
	return response.Success(c, nil, "setup completed")
}

// User handlers
func (s *Server) handleGetMe(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	user := middleware.GetUser(c)
	if user == nil {
		return response.Unauthorized(c, "authentication required")
	}

	userDetails, err := s.container.UserService.GetByID(ctx, user.ID)
	if err != nil {
		return response.InternalError(c, "failed to get user details")
	}
	return response.Success(c, userDetails.ToResponse(), "user profile")
}

func (s *Server) handleGetProfile(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	userID := c.Params("id")
	if userID == "" {
		return response.BadRequest(c, "user ID is required", nil)
	}

	userDetails, err := s.container.UserService.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return response.NotFound(c, "user not found")
		}
		return response.InternalError(c, "failed to get user profile")
	}
	return response.Success(c, userDetails.ToResponse(), "user profile")
}

func (s *Server) handleUpdateProfile(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	userID := c.Params("id")
	if userID == "" {
		return response.BadRequest(c, "user ID is required", nil)
	}

	currentUser := middleware.GetUser(c)
	if currentUser == nil {
		return response.Unauthorized(c, "authentication required")
	}

	var req user.UpdateRequest
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
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

	updatedUser, err := s.container.UserService.Update(ctx, userID, req, domainUser)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return response.NotFound(c, "user not found")
		}
		if errors.Is(err, user.ErrInsufficientPermissions) {
			return response.Forbidden(c, "insufficient permissions")
		}
		return response.InternalError(c, "failed to update user")
	}
	return response.Success(c, updatedUser.ToResponse(), "user updated")
}

func (s *Server) handleListUsers(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, pagination, err := s.container.UserService.List(ctx, page, limit)
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

func (s *Server) handleCreateUser(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
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
		if errors.Is(err, user.ErrEmailAlreadyExists) {
			return response.Conflict(c, "email already exists")
		}
		return response.InternalError(c, "failed to create user")
	}
	return response.Success(c, newUser.ToResponse(), "user created")
}

// Source handlers
func (s *Server) handleCreateSource(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	currentUser := middleware.GetUser(c)
	if currentUser == nil {
		return response.Unauthorized(c, "authentication required")
	}

	var req source.CreateRequest
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	// Convert middleware User to domain CurrentUser
	domainUser := &source.CurrentUser{
		ID:   currentUser.ID,
		Role: string(currentUser.Role),
	}

	newSource, err := s.container.SourceService.Create(ctx, req, domainUser)
	if err != nil {
		if errors.Is(err, source.ErrSourceAlreadyExists) {
			return response.Conflict(c, "source already exists")
		}
		return response.InternalError(c, "failed to create source")
	}
	sourceResp, err := newSource.ToResponse()
	if err != nil {
		return response.InternalError(c, "failed to format source response")
	}
	return response.Success(c, sourceResp, "source created")
}

func (s *Server) handleGetSource(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	sourceID := c.Params("id")
	if sourceID == "" {
		return response.BadRequest(c, "source ID is required", nil)
	}

	sourceData, err := s.container.SourceService.GetByID(ctx, sourceID)
	if err != nil {
		if errors.Is(err, source.ErrSourceNotFound) {
			return response.NotFound(c, "source not found")
		}
		return response.InternalError(c, "failed to get source")
	}
	sourceResp, err := sourceData.ToResponse()
	if err != nil {
		return response.InternalError(c, "failed to format source response")
	}
	return response.Success(c, sourceResp, "source details")
}

func (s *Server) handleListSources(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	sources, pagination, err := s.container.SourceService.List(ctx, page, limit)
	if err != nil {
		return response.InternalError(c, "failed to list sources")
	}

	listResp := &source.ListSourcesResponse{
		Sources:    sources,
		Pagination: pagination,
	}

	return response.Success(c, listResp, "sources listed")
}

func (s *Server) handleUpdateSource(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	sourceID := c.Params("id")
	if sourceID == "" {
		return response.BadRequest(c, "source ID is required", nil)
	}

	currentUser := middleware.GetUser(c)
	if currentUser == nil {
		return response.Unauthorized(c, "authentication required")
	}

	var req source.UpdateRequest
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	// Convert middleware User to domain CurrentUser
	domainUser := &source.CurrentUser{
		ID:   currentUser.ID,
		Role: string(currentUser.Role),
	}

	updatedSource, err := s.container.SourceService.Update(ctx, sourceID, req, domainUser)
	if err != nil {
		if errors.Is(err, source.ErrSourceNotFound) {
			return response.NotFound(c, "source not found")
		}
		if errors.Is(err, source.ErrSourceAlreadyExists) {
			return response.Conflict(c, "source name already exists")
		}
		return response.InternalError(c, "failed to update source")
	}
	sourceResp, err := updatedSource.ToResponse()
	if err != nil {
		return response.InternalError(c, "failed to format source response")
	}
	return response.Success(c, sourceResp, "source updated")
}

func (s *Server) handleDeleteSource(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	sourceID := c.Params("id")
	if sourceID == "" {
		return response.BadRequest(c, "source ID is required", nil)
	}

	err := s.container.SourceService.Delete(ctx, sourceID)
	if err != nil {
		if errors.Is(err, source.ErrSourceNotFound) {
			return response.NotFound(c, "source not found")
		}
		return response.InternalError(c, "failed to delete source")
	}
	return response.Success(c, nil, "source deleted")
}
