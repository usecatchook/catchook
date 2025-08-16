package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	destination "github.com/theotruvelot/catchook/internal/destination/domain"
	"github.com/theotruvelot/catchook/internal/platform/auth"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
)

type destinationService struct {
	destinationRepo destination.Repository
	appLogger       logger.Logger
}

func NewDestinationService(destinationRepo destination.Repository, appLogger logger.Logger) destination.Service {
	return &destinationService{
		destinationRepo: destinationRepo,
		appLogger:       appLogger,
	}
}

func (s destinationService) Create(ctx context.Context, req destination.CreateRequest) (*destination.Destination, error) {
	ctx, span := tracer.StartSpan(ctx, "destination.service.create")
	defer span.End()

	s.appLogger.Info(ctx, "Creating new destination", logger.String("name", req.Name))

	exist, err := s.destinationRepo.GetByName(ctx, req.Name)
	if err != nil {
		span.RecordError(err)
		s.appLogger.Error(ctx, "Failed to check existing destination by name", logger.Error(err))
		return nil, fmt.Errorf("checking existing destination by name: %w", err)
	}
	if exist != nil {
		return nil, destination.ErrDestinationAlreadyExists
	}

	config, err := validateAndMarshalConfig(req.DestinationType, req.Config)
	if err != nil {
		span.RecordError(err)
		s.appLogger.Error(ctx, "Failed to build config", logger.Error(err))
		return nil, fmt.Errorf("building config: %w", err)
	}

	currentUserID, err := auth.GetUserID(ctx)
	if err != nil {
		span.RecordError(err)
		s.appLogger.Error(ctx, "Failed to get current user id", logger.Error(err))
		return nil, fmt.Errorf("getting current user id: %w", err)
	}

	newDestination := &destination.Destination{
		UserID:          currentUserID,
		Name:            req.Name,
		Description:     req.Description,
		DestinationType: req.DestinationType,
		Config:          config,
		IsActive:        true,
		DelaySeconds:    req.DelaySeconds,
		RetryAttempts:   req.RetryAttempts,
	}

	if err := s.destinationRepo.Create(ctx, newDestination); err != nil {
		span.RecordError(err)
		s.appLogger.Error(ctx, "Failed to create destination", logger.Error(err))
		return nil, fmt.Errorf("creating destination: %w", err)
	}

	return newDestination, nil
}

func validateAndMarshalConfig(destType destination.DestinationType, cfg map[string]interface{}) (string, error) {
	errors := map[string]string{}

	switch destType {
	case destination.DestinationTypeHTTP:
		httpConfig, err := destination.HTTPConfigFromMap(cfg)
		if err != nil {
			errors["config"] = fmt.Sprintf("invalid HTTP config format: %v", err)
		} else {
			if err := httpConfig.Validate(); err != nil {
				errors["config"] = fmt.Sprintf("HTTP config validation failed: %v", err)
			} else {
				cfg, err = httpConfig.ToMap()
				if err != nil {
					errors["config"] = fmt.Sprintf("failed to convert HTTP config: %v", err)
				}
			}
		}
	case destination.DestinationTypeRabbitMQ:
		requireFields(errors, cfg, "host", "queue")
	case destination.DestinationTypeDatabase:
		requireFields(errors, cfg, "connection_string", "table")
	case destination.DestinationTypeFile:
		requireFields(errors, cfg, "path")
	case destination.DestinationTypeQueue:
		requireFields(errors, cfg, "host", "queue")
	case destination.DestinationTypeCLI:
		requireFields(errors, cfg, "command")
	default:
		return "", &validatorpkg.ValidationErrors{Errors: map[string]string{
			"destination_type": "unsupported destination_type",
		}}
	}

	if len(errors) > 0 {
		return "", &validatorpkg.ValidationErrors{Errors: errors}
	}

	if cfg == nil {
		cfg = map[string]interface{}{}
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal config: %w", err)
	}
	return string(b), nil
}

func requireFields(errs map[string]string, cfg map[string]interface{}, fields ...string) {
	for _, f := range fields {
		v, ok := cfg[f]
		if !ok {
			errs["config."+f] = "is required"
			continue
		}
		switch val := v.(type) {
		case string:
			if strings.TrimSpace(val) == "" {
				errs["config."+f] = "cannot be empty"
			}
		}
	}
}

func (s destinationService) GetByID(ctx context.Context, id string) (*destination.Destination, error) {
	ctx, span := tracer.StartSpan(ctx, "destination.service.get_by_id")
	defer span.End()

	s.appLogger.Info(ctx, "Fetching destination by ID", logger.String("destination_id", id))

	destinationData, err := s.destinationRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("getting destination by ID: %w", err)
	}
	if destinationData == nil {
		return nil, destination.ErrDestinationNotFound
	}

	return destinationData, nil
}

func (s destinationService) List(ctx context.Context, req destination.ListDestinationsRequest) ([]*destination.DestinationListItem, *response.Pagination, error) {
	ctx, span := tracer.StartSpan(ctx, "destination.service.list")
	defer span.End()

	destinations, pagination, err := s.destinationRepo.List(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.appLogger.Error(ctx, "Failed to list destinations", logger.Error(err))
		return nil, nil, fmt.Errorf("failed to list destinations: %w", err)
	}

	return destinations, pagination, nil
}

func (s destinationService) Update(ctx context.Context, id string, req destination.UpdateRequest) (*destination.Destination, error) {
	ctx, span := tracer.StartSpan(ctx, "destination.service.update")
	defer span.End()

	s.appLogger.Info(ctx, "Updating destination", logger.String("destination_id", id))
	existing, err := s.destinationRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("getting destination by ID: %w", err)
	}
	if existing == nil {
		return nil, destination.ErrDestinationNotFound
	}

	var configStr string
	if req.Config != nil {
		destType := req.DestinationType
		if destType == "" {
			destType = existing.DestinationType
		}
		configStr, err = validateAndMarshalConfig(destType, req.Config)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("building config: %w", err)
		}
	}
	updated, err := s.destinationRepo.Update(ctx, id, req.Name, req.Description, req.DestinationType, configStr, req.IsActive, req.DelaySeconds, req.RetryAttempts)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("updating destination: %w", err)
	}

	return updated, nil
}

func (s destinationService) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.StartSpan(ctx, "destination.service.delete")
	defer span.End()

	s.appLogger.Info(ctx, "Deleting destination", logger.String("destination_id", id))

	existing, err := s.destinationRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.appLogger.Error(ctx, "Failed to get destination by ID", logger.Error(err))
		return fmt.Errorf("getting destination by ID: %w", err)
	}
	if existing == nil {
		return destination.ErrDestinationNotFound
	}

	if err := s.destinationRepo.Delete(ctx, id); err != nil {
		span.RecordError(err)
		s.appLogger.Error(ctx, "Failed to delete destination", logger.Error(err))
		return fmt.Errorf("deleting destination: %w", err)
	}

	return nil
}
