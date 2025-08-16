package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	destination "github.com/theotruvelot/catchook/internal/destination/domain"
	"github.com/theotruvelot/catchook/internal/platform/storage/postgres/generated"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
)

type destinationRepository struct {
	db        *pgxpool.Pool
	queries   *generated.Queries
	appLogger logger.Logger
}

func NewDestinationRepository(db *pgxpool.Pool, appLogger logger.Logger) destination.Repository {
	return &destinationRepository{
		db:        db,
		queries:   generated.New(db),
		appLogger: appLogger,
	}
}

func (r destinationRepository) Create(ctx context.Context, dest *destination.Destination) error {
	ctx, span := tracer.StartSpan(ctx, "destination.repository.create")
	defer span.End()

	userId, err := uuid.Parse(dest.UserID)
	if err != nil {
		span.RecordError(err)
		r.appLogger.Error(ctx, "Invalid user id", logger.Error(err))
		return fmt.Errorf("invalid user id: %w", err)
	}

	config := dest.Config
	if config == "" {
		config = "{}"
	}

	result, err := r.queries.CreateDestination(ctx,
		userId,
		dest.Name,
		dest.Description,
		generated.DestinationType(dest.DestinationType),
		[]byte(config),
		dest.IsActive,
		dest.DelaySeconds,
		dest.RetryAttempts,
	)

	if err != nil {
		r.appLogger.Error(ctx, "Failed to create destination",
			logger.String("name", dest.Name),
			logger.String("user_id", dest.UserID),
			logger.Error(err),
		)
		span.RecordError(err)
		return fmt.Errorf("failed to create destination: %w", err)
	}

	dest.ID = result.ID.String()
	dest.CreatedAt = result.CreatedAt.Time
	dest.UpdatedAt = result.UpdatedAt.Time

	return nil
}

func (r destinationRepository) GetByID(ctx context.Context, id string) (*destination.Destination, error) {
	ctx, span := tracer.StartSpan(ctx, "destination.repository.get_by_id")
	defer span.End()

	uid, err := uuid.Parse(id)
	if err != nil {
		r.appLogger.Error(ctx, "Invalid UUID format",
			logger.String("destination_id", id),
			logger.Error(err),
		)
		return nil, fmt.Errorf("invalid destination ID format: %w", err)
	}

	result, err := r.queries.GetDestinationByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		r.appLogger.Error(ctx, "Failed to get destination by ID", logger.Error(err))
		return nil, fmt.Errorf("failed to get destination by ID: %w", err)
	}

	return &destination.Destination{
		ID:              result.ID.String(),
		UserID:          result.UserID.String(),
		Name:            result.Name,
		Description:     result.Description,
		DestinationType: destination.DestinationType(result.DestinationType),
		Config:          string(result.Config),
		IsActive:        result.IsActive,
		DelaySeconds:    result.DelaySeconds,
		RetryAttempts:   result.RetryAttempts,
		CreatedAt:       result.CreatedAt.Time,
		UpdatedAt:       result.UpdatedAt.Time,
	}, nil
}

func (r destinationRepository) GetByName(ctx context.Context, name string) (*destination.Destination, error) {
	ctx, span := tracer.StartSpan(ctx, "destination.repository.get_by_name")
	defer span.End()

	result, err := r.queries.GetDestinationByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		r.appLogger.Error(ctx, "Failed to get destination by name", logger.Error(err))
		span.RecordError(err)
		return nil, err
	}

	return &destination.Destination{
		ID:              result.ID.String(),
		UserID:          result.UserID.String(),
		Name:            result.Name,
		Description:     result.Description,
		DestinationType: destination.DestinationType(result.DestinationType),
		Config:          string(result.Config),
		IsActive:        result.IsActive,
		DelaySeconds:    result.DelaySeconds,
		RetryAttempts:   result.RetryAttempts,
		CreatedAt:       result.CreatedAt.Time,
		UpdatedAt:       result.UpdatedAt.Time,
	}, nil
}

func (r destinationRepository) List(ctx context.Context, req destination.ListDestinationsRequest) ([]*destination.DestinationListItem, *response.Pagination, error) {
	ctx, span := tracer.StartSpan(ctx, "destination.repository.list")
	defer span.End()

	offset := (req.Page - 1) * req.Limit
	filterByActive := req.IsActive != nil
	isActiveValue := false
	if req.IsActive != nil {
		isActiveValue = *req.IsActive
	}

	results, err := r.queries.ListDestinations(ctx,
		req.Search,
		req.DestinationType,
		filterByActive,
		isActiveValue,
		req.OrderBy,
		req.Order,
		int32(req.Limit),
		int32(offset),
	)
	if err != nil {
		span.RecordError(err)
		return nil, nil, fmt.Errorf("failed to list destinations: %w", err)
	}

	total, err := r.queries.CountDestinations(ctx, req.Search, req.DestinationType, filterByActive, isActiveValue)
	if err != nil {
		span.RecordError(err)
		return nil, nil, fmt.Errorf("failed to count destinations: %w", err)
	}

	destinations := make([]*destination.DestinationListItem, len(results))
	for i, result := range results {
		destinations[i] = &destination.DestinationListItem{
			Name:            result.Name,
			Description:     result.Description,
			DestinationType: string(result.DestinationType),
			IsActive:        result.IsActive,
			CreatedAt:       result.CreatedAt.Time,
			UpdatedAt:       result.UpdatedAt.Time,
		}
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	if totalPages < 1 {
		totalPages = 1
	}

	pagination := &response.Pagination{
		CurrentPage: req.Page,
		TotalPages:  totalPages,
		Total:       int(total),
		Limit:       req.Limit,
		HasNext:     req.Page < totalPages,
		HasPrev:     req.Page > 1,
	}

	return destinations, pagination, nil
}

func (r destinationRepository) Update(ctx context.Context, id, name, description string, destType destination.DestinationType, config string, isActive bool, delaySeconds, retryAttempts int32) (*destination.Destination, error) {
	ctx, span := tracer.StartSpan(ctx, "destination.repository.update")
	defer span.End()

	uid, err := uuid.Parse(id)
	if err != nil {
		span.RecordError(err)
		r.appLogger.Error(ctx, "Invalid destination id", logger.Error(err))
		return nil, fmt.Errorf("invalid destination id: %w", err)
	}

	configParam := []byte(config)
	if config == "" {
		configParam = nil
	}

	result, err := r.queries.UpdateDestination(ctx,
		uid,
		name,
		description,
		generated.DestinationType(destType),
		configParam,
		isActive,
		delaySeconds,
		retryAttempts,
	)
	if err != nil {
		span.RecordError(err)
		r.appLogger.Error(ctx, "Failed to update destination", logger.Error(err))
		return nil, fmt.Errorf("failed to update destination: %w", err)
	}

	return &destination.Destination{
		ID:              result.ID.String(),
		UserID:          result.UserID.String(),
		Name:            result.Name,
		Description:     result.Description,
		DestinationType: destination.DestinationType(result.DestinationType),
		Config:          string(result.Config),
		IsActive:        result.IsActive,
		DelaySeconds:    result.DelaySeconds,
		RetryAttempts:   result.RetryAttempts,
		CreatedAt:       result.CreatedAt.Time,
		UpdatedAt:       result.UpdatedAt.Time,
	}, nil
}

func (r destinationRepository) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.StartSpan(ctx, "destination.repository.delete")
	defer span.End()

	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid destination id: %w", err)
	}

	if err := r.queries.DeleteDestination(ctx, uid); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to delete destination: %w", err)
	}
	return nil
}
