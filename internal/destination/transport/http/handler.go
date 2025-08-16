package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	destination "github.com/theotruvelot/catchook/internal/destination/domain"
	"github.com/theotruvelot/catchook/internal/platform/http/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
)

type Handler struct {
	destinationService destination.Service
	validator          *validatorpkg.Validator
}

func NewHandler(destinationService destination.Service, validator *validatorpkg.Validator) *Handler {
	return &Handler{
		destinationService: destinationService,
		validator:          validator,
	}
}

func (h *Handler) CreateDestination(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)

	ctx, span := tracer.StartSpan(ctx, "destination.handler.create")
	defer span.End()

	var req destination.CreateRequest
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	destinationResp, err := h.destinationService.Create(ctx, req)
	if err != nil {
		var verr *validatorpkg.ValidationErrors
		switch {
		case errors.As(err, &verr):
			return response.ValidationFailed(c, verr.Errors)
		case errors.Is(err, destination.ErrDestinationAlreadyExists):
			return response.Conflict(c, "destination already exists")
		default:
			return response.InternalError(c, "failed to create destination")
		}
	}

	resp, err := destinationResp.ToResponse()
	if err != nil {
		return response.InternalError(c, "failed to serialize destination")
	}

	return response.Success(c, resp, "destination created")
}

func (h *Handler) GetDestination(c *fiber.Ctx) error {
	ctx, span := tracer.StartSpan(middleware.GetContextWithRequestID(c), "destination.handler.get")
	defer span.End()

	destinationID := c.Params("id")
	if destinationID == "" {
		return response.BadRequest(c, "destination_id is required", nil)
	}

	destinationResp, err := h.destinationService.GetByID(ctx, destinationID)
	if err != nil {
		if errors.Is(err, destination.ErrDestinationNotFound) {
			return response.NotFound(c, "destination not found")
		}
		return response.InternalError(c, "failed to get destination")
	}

	resp, err := destinationResp.ToResponse()
	if err != nil {
		return response.InternalError(c, "failed to serialize destination")
	}
	return response.Success(c, resp, "destination")
}

func (h *Handler) ListDestinations(c *fiber.Ctx) error {
	ctx, span := tracer.StartSpan(middleware.GetContextWithRequestID(c), "destination.handler.list")
	defer span.End()

	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	limit := c.QueryInt("limit", 20)
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	orderBy := c.Query("order_by", "created_at")
	validOrderBys := map[string]bool{
		"name":       true,
		"created_at": true,
		"updated_at": true,
		"is_active":  true,
	}
	if !validOrderBys[orderBy] {
		return response.BadRequest(c, "invalid order_by. Must be one of: name, created_at, updated_at", nil)
	}

	order := c.Query("order", "desc")
	if order != "asc" && order != "desc" {
		return response.BadRequest(c, "invalid order. Must be 'asc' or 'desc'", nil)
	}

	destType := c.Query("destination_type")
	if destType != "" {
		validTypes := map[string]bool{
			"http":     true,
			"rabbitmq": true,
			"database": true,
			"file":     true,
			"queue":    true,
			"cli":      true,
		}
		if !validTypes[destType] {
			return response.BadRequest(c, "invalid destination_type", nil)
		}
	}

	req := destination.ListDestinationsRequest{
		Page:            page,
		Limit:           limit,
		DestinationType: destType,
		Search:          c.Query("search"),
		OrderBy:         orderBy,
		Order:           order,
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		switch isActiveStr {
		case "true":
			req.IsActive = &[]bool{true}[0]
		case "false":
			req.IsActive = &[]bool{false}[0]
		default:
			return response.BadRequest(c, "invalid is_active. Must be 'true' or 'false'", nil)
		}
	}

	destinations, pagination, err := h.destinationService.List(ctx, req)
	if err != nil {
		return response.InternalError(c, "failed to list destinations")
	}

	return response.Success(c, &destination.ListDestinationsResponse{
		Destinations: destinations,
		Pagination:   pagination,
	}, "destinations listed")
}

func (h *Handler) UpdateDestination(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	ctx, span := tracer.StartSpan(ctx, "destination.handler.update")
	defer span.End()

	destinationID := c.Params("id")
	if destinationID == "" {
		return response.BadRequest(c, "destination_id is required", nil)
	}

	var req destination.UpdateRequest
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		var verr *validatorpkg.ValidationErrors
		if errors.As(err, &verr) {
			return response.ValidationFailed(c, verr.Errors)
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	updated, err := h.destinationService.Update(ctx, destinationID, req)
	if err != nil {
		var verr *validatorpkg.ValidationErrors
		switch {
		case errors.As(err, &verr):
			return response.ValidationFailed(c, verr.Errors)
		case errors.Is(err, destination.ErrDestinationNotFound):
			return response.NotFound(c, "destination not found")
		case errors.Is(err, destination.ErrDestinationAlreadyExists):
			return response.Conflict(c, "destination already exists")
		default:
			return response.InternalError(c, "failed to update destination")
		}
	}

	resp, err := updated.ToResponse()
	if err != nil {
		return response.InternalError(c, "failed to serialize destination")
	}

	return response.Success(c, resp, "destination updated")
}

func (h *Handler) DeleteDestination(c *fiber.Ctx) error {
	ctx := middleware.GetContextWithRequestID(c)
	ctx, span := tracer.StartSpan(ctx, "destination.handler.delete")
	defer span.End()

	destinationID := c.Params("id")
	if destinationID == "" {
		return response.BadRequest(c, "destination_id is required", nil)
	}

	if err := h.destinationService.Delete(ctx, destinationID); err != nil {
		switch {
		case errors.Is(err, destination.ErrDestinationNotFound):
			return response.NotFound(c, "destination not found")
		default:
			return response.InternalError(c, "failed to delete destination")
		}
	}

	return response.Success(c, nil, "destination deleted")
}
