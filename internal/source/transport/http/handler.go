package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/theotruvelot/catchook/internal/platform/http/middleware"
	source "github.com/theotruvelot/catchook/internal/source/domain"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
)

// Handler holds the source-specific dependencies
type Handler struct {
	sourceService source.Service
	validator     *validatorpkg.Validator
}

// NewHandler creates a new source handler
func NewHandler(sourceService source.Service, validator *validatorpkg.Validator) *Handler {
	return &Handler{
		sourceService: sourceService,
		validator:     validator,
	}
}

func (h *Handler) CreateSource(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	ctx, span := tracer.StartSpan(ctx, "source.handler")
	defer span.End()

	var req source.CreateRequest
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	sourceResp, err := h.sourceService.Create(ctx, req)
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

func (h *Handler) GetSource(c *fiber.Ctx) error {
	ctx, span := tracer.StartSpan(middleware.GetContextWithRequestID(c), "source.handler.get")
	defer span.End()

	sourceID := c.Params("id")
	if sourceID == "" {
		return response.BadRequest(c, "source_id is required", nil)
	}
	sourceResp, err := h.sourceService.GetByID(ctx, sourceID)
	if err != nil {
		if errors.Is(err, source.ErrSourceNotFound) {
			return response.NotFound(c, "source not found")
		}
		return response.InternalError(c, "failed to get source")
	}

	resp, err := sourceResp.ToResponse()
	if err != nil {
		return response.InternalError(c, "failed to serialize source")
	}
	return response.Success(c, resp, "source")
}

func (h *Handler) ListSources(c *fiber.Ctx) error {
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

	sources, pagination, err := h.sourceService.List(ctx, page, limit)
	if err != nil {
		return response.InternalError(c, "failed to list sources")
	}

	listResp := &source.ListSourcesResponse{
		Sources:    sources,
		Pagination: pagination,
	}

	return response.Success(c, listResp, "sources listed")
}

func (h *Handler) UpdateSource(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	ctx, span := tracer.StartSpan(ctx, "source.handler.update")
	defer span.End()

	// L'auth est maintenant gérée par le middleware RequirePermission
	// Plus besoin de vérifier manuellement l'authentification

	sourceID := c.Params("id")
	if sourceID == "" {
		return response.BadRequest(c, "source_id is required", nil)
	}

	var req source.UpdateRequest
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	updated, err := h.sourceService.Update(ctx, sourceID, req)
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

func (h *Handler) DeleteSource(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	ctx, span := tracer.StartSpan(ctx, "source.handler.delete")
	defer span.End()

	// L'auth est maintenant gérée par le middleware RequirePermission
	// Plus besoin de vérifier manuellement l'authentification

	sourceID := c.Params("id")
	if sourceID == "" {
		return response.BadRequest(c, "source_id is required", nil)
	}

	if err := h.sourceService.Delete(ctx, sourceID); err != nil {
		switch {
		case errors.Is(err, source.ErrSourceNotFound):
			return response.NotFound(c, "source not found")
		default:
			return response.InternalError(c, "failed to delete source")
		}
	}

	return response.Success(c, nil, "source deleted")
}
