package server

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/domain/source"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
	"github.com/theotruvelot/catchook/storage/postgres/generated"
)

func (s *Server) handleCreateSource(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	ctx, span := tracer.StartSpan(ctx, "source.handler")
	defer span.End()
	currentUser := middleware.GetUser(c)
	if currentUser == nil || currentUser.Role == generated.UserRoleViewer {
		return response.Forbidden(c, "viewer cannot create sources")
	}

	var req source.CreateRequest
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	sourceResp, err := s.container.SourceService.Create(ctx, req, currentUser)
	if err != nil {
		var verr *validatorpkg.ValidationErrors
		switch {
		case errors.As(err, &verr):
			return response.ValidationFailed(c, verr.Errors)
		case errors.Is(err, source.ErrSourceAlreadyExists):
			return response.Conflict(c, "source already exists")
		default:
			return response.InternalError(c, "failed to create source")
		}
	}

	resp, err := sourceResp.ToResponse()
	if err != nil {
		return response.InternalError(c, "failed to serialize source")
	}

	return response.Success(c, resp, "source created")
}

func (s *Server) handleGetSource(c *fiber.Ctx) error {
	ctx, span := tracer.StartSpan(middleware.GetContextWithRequestID(c), "source.handler.get")
	defer span.End()

	sourceID := c.Params("id")
	if sourceID == "" {
		return response.BadRequest(c, "source_id is required", nil)
	}
	sourceResp, err := s.container.SourceService.GetByID(ctx, sourceID)
	if err != nil {
		return response.NotFound(c, "source not found")
	}

	resp, err := sourceResp.ToResponse()
	if err != nil {
		return response.InternalError(c, "failed to serialize source")
	}
	return response.Success(c, resp, "source")
}

func (s *Server) handleListSources(c *fiber.Ctx) error {
	ctx, span := tracer.StartSpan(middleware.GetContextWithRequestID(c), "source.handler.list")
	defer span.End()

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

	sources, pagination, err := s.container.SourceService.List(ctx, page, limit)
	if err != nil {
		return response.InternalError(c, "failed to list sources")
	}

	return response.Paginated(c, sources, *pagination, "sources list")
}

func (s *Server) handleUpdateSource(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	ctx, span := tracer.StartSpan(ctx, "source.handler.update")
	defer span.End()

	currentUser := middleware.GetUser(c)
	if currentUser == nil || currentUser.Role == generated.UserRoleViewer {
		return response.Forbidden(c, "viewer cannot update sources")
	}

	sourceID := c.Params("id")
	if sourceID == "" {
		return response.BadRequest(c, "source_id is required", nil)
	}

	var req source.UpdateRequest
	if err := s.container.Validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	updated, err := s.container.SourceService.Update(ctx, sourceID, req, currentUser)
	if err != nil {
		var verr *validatorpkg.ValidationErrors
		switch {
		case errors.As(err, &verr):
			return response.ValidationFailed(c, verr.Errors)
		case errors.Is(err, source.ErrSourceNotFound):
			return response.NotFound(c, "source not found")
		case errors.Is(err, source.ErrSourceAlreadyExists):
			return response.Conflict(c, "source already exists")
		default:
			return response.InternalError(c, "failed to update source")
		}
	}

	resp, err := updated.ToResponse()
	if err != nil {
		return response.InternalError(c, "failed to serialize source")
	}

	return response.Success(c, resp, "source updated")
}

func (s *Server) handleDeleteSource(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	ctx, span := tracer.StartSpan(ctx, "source.handler.delete")
	defer span.End()

	currentUser := middleware.GetUser(c)
	if currentUser == nil || currentUser.Role == generated.UserRoleViewer {
		return response.Forbidden(c, "viewer cannot delete sources")
	}

	sourceID := c.Params("id")
	if sourceID == "" {
		return response.BadRequest(c, "source_id is required", nil)
	}

	if err := s.container.SourceService.Delete(ctx, sourceID); err != nil {
		switch {
		case errors.Is(err, source.ErrSourceNotFound):
			return response.NotFound(c, "source not found")
		default:
			return response.InternalError(c, "failed to delete source")
		}
	}

	return response.NoContent(c)
}
