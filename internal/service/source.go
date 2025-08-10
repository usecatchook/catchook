package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/theotruvelot/catchook/internal/domain/source"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
	validatorpkg "github.com/theotruvelot/catchook/pkg/validator"
)

type sourceService struct {
	sourceRepo source.Repository
	appLogger  logger.Logger
}

func NewSourceService(sourceRepo source.Repository, appLogger logger.Logger) source.Service {
	return &sourceService{
		sourceRepo: sourceRepo,
		appLogger:  appLogger,
	}
}

func (s sourceService) Create(ctx context.Context, req source.CreateRequest, currentUser *middleware.User) (*source.Source, error) {
	ctx, span := tracer.StartSpan(ctx, "source.service.create")
	defer span.End()

	s.appLogger.Info(ctx, "Creating new source", logger.String("name", req.Name))

	exist, err := s.sourceRepo.GetByName(ctx, req.Name)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("checking existing source by name: %w", err)
	}
	if exist != nil {
		return nil, source.ErrSourceAlreadyExists
	}

	authCfg, err := validateAndMarshalAuthConfig(req.AuthType, req.AuthConfig)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("building auth config: %w", err)
	}

	newSource := &source.Source{
		UserID:      currentUser.ID,
		Name:        req.Name,
		Description: req.Description,
		Protocol:    req.Protocol,
		AuthType:    req.AuthType,
		AuthConfig:  authCfg,
		IsActive:    true,
	}

	if err := s.sourceRepo.Create(ctx, newSource); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("creating source: %w", err)
	}

	return newSource, nil
}

func validateAndMarshalAuthConfig(authType source.AuthType, cfg map[string]any) (string, error) {
	errors := map[string]string{}

	switch authType {
	case source.AuthTypeNone:
		return "{}", nil
	case source.AuthTypeBasic:
		requireFields(errors, cfg, "username", "password")
	case source.AuthTypeBearer:
		requireFields(errors, cfg, "token")
	case source.AuthTypeApikey:
		requireFields(errors, cfg, "location", "value")
	case source.AuthTypeSignature:
		requireFields(errors, cfg, "secret", "header", "algorithm", "encoding")
		if algo, ok := getString(cfg, "algorithm"); ok {
			switch algo {
			case "sha-1", "sha-256", "sha-512", "md5":
			default:
				errors["auth_config.algorithm"] = "must be one of: sha-1 sha-256 sha-512 md5"
			}
		}
		if enc, ok := getString(cfg, "encoding"); ok {
			switch enc {
			case "base64", "base64url", "hex":
			default:
				errors["auth_config.encoding"] = "must be one of: base64 base64url hex"
			}
		}
	default:
		return "", &validatorpkg.ValidationErrors{Errors: map[string]string{
			"auth_type": "unsupported auth_type",
		}}
	}

	if len(errors) > 0 {
		return "", &validatorpkg.ValidationErrors{Errors: errors}
	}

	if cfg == nil {
		cfg = map[string]any{}
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal auth_config: %w", err)
	}
	return string(b), nil
}

func requireFields(errs map[string]string, cfg map[string]any, fields ...string) {
	for _, f := range fields {
		v, ok := cfg[f]
		if !ok {
			errs["auth_config."+f] = "is required"
			continue
		}
		switch val := v.(type) {
		case string:
			if strings.TrimSpace(val) == "" {
				errs["auth_config."+f] = "cannot be empty"
			}
		}
	}
}

func getString(cfg map[string]any, key string) (string, bool) {
	if cfg == nil {
		return "", false
	}
	v, ok := cfg[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	return s, true
}

func (s sourceService) GetByID(ctx context.Context, id string) (*source.Source, error) {
	ctx, span := tracer.StartSpan(ctx, "source.service.get_by_id")
	defer span.End()

	s.appLogger.Info(ctx, "Fetching source by ID", logger.String("source_id", id))

	sourceData, err := s.sourceRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("getting source by ID: %w", err)
	}
	if sourceData == nil {
		return nil, source.ErrSourceNotFound
	}

	return sourceData, nil
}

func (s sourceService) List(ctx context.Context, page, limit int) ([]*source.SourceResponse, *response.Pagination, error) {
	ctx, span := tracer.StartSpan(ctx, "source.service.list")
	defer span.End()

	sources, meta, err := s.sourceRepo.List(ctx, page, limit)
	if err != nil {
		span.RecordError(err)
		s.appLogger.Error(ctx, "Failed to list sources", logger.Error(err))
		return nil, nil, fmt.Errorf("failed to list sources: %w", err)
	}

	respList, err := source.ToResponses(sources)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize source: %w", err)
	}

	return respList, meta, nil
}

func (s sourceService) Update(ctx context.Context, id string, req source.UpdateRequest, currentUser *middleware.User) (*source.Source, error) {
	ctx, span := tracer.StartSpan(ctx, "source.service.update")
	defer span.End()

	s.appLogger.Info(ctx, "Updating source", logger.String("source_id", id))

	existing, err := s.sourceRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("getting source by ID: %w", err)
	}
	if existing == nil {
		return nil, source.ErrSourceNotFound
	}

	if req.Name != "" && req.Name != existing.Name {
		other, err := s.sourceRepo.GetByName(ctx, req.Name)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("checking existing source by name: %w", err)
		}
		if other != nil && other.ID != id {
			return nil, source.ErrSourceAlreadyExists
		}
	}

	finalAuthType := existing.AuthType
	if req.AuthType != "" {
		finalAuthType = req.AuthType
	}

	var finalAuthConfig string
	switch {
	case req.AuthConfig != nil:
		finalAuthConfig, err = validateAndMarshalAuthConfig(finalAuthType, req.AuthConfig)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("building auth config: %w", err)
		}
	case req.AuthType != "":
		return nil, fmt.Errorf("auth_config is required when changing auth_type")
	default:
		finalAuthConfig = existing.AuthConfig
	}

	name := existing.Name
	if strings.TrimSpace(req.Name) != "" {
		name = req.Name
	}
	description := existing.Description
	if req.Description != "" {
		description = req.Description
	}
	protocol := existing.Protocol
	if strings.TrimSpace(req.Protocol) != "" {
		protocol = req.Protocol
	}

	updated := &source.Source{
		ID:          existing.ID,
		UserID:      existing.UserID,
		Name:        name,
		Description: description,
		Protocol:    protocol,
		AuthType:    finalAuthType,
		AuthConfig:  finalAuthConfig,
		IsActive:    existing.IsActive,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   existing.UpdatedAt,
	}

	if err := s.sourceRepo.Update(ctx, updated); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("updating source: %w", err)
	}

	return updated, nil
}

func (s sourceService) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.StartSpan(ctx, "source.service.delete")
	defer span.End()

	s.appLogger.Info(ctx, "Deleting source", logger.String("source_id", id))

	existing, err := s.sourceRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("getting source by ID: %w", err)
	}
	if existing == nil {
		return source.ErrSourceNotFound
	}

	if err := s.sourceRepo.Delete(ctx, id); err != nil {
		span.RecordError(err)
		return fmt.Errorf("deleting source: %w", err)
	}

	return nil
}
