package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/theotruvelot/catchook/internal/domain/source"
	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
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

	authCfg, err := buildAuthConfigJSON(req)
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

func buildAuthConfigJSON(req source.CreateRequest) (string, error) {
	switch req.AuthType {
	case source.AuthTypeNone:
		return "{}", nil

	case source.AuthTypeBasic:
		if req.BasicAuth == nil {
			return "", fmt.Errorf("basic_auth is required when auth_type=basic")
		}
		b, err := json.Marshal(struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: req.BasicAuth.Username,
			Password: req.BasicAuth.Password,
		})
		if err != nil {
			return "", err
		}
		return string(b), nil

	case source.AuthTypeBearer:
		if req.BearerAuth == nil {
			return "", fmt.Errorf("bearer_auth is required when auth_type=bearer")
		}
		b, err := json.Marshal(struct {
			Token string `json:"token"`
		}{
			Token: req.BearerAuth.Token,
		})
		if err != nil {
			return "", err
		}
		return string(b), nil

	case source.AuthTypeApikey:
		if req.APIKeyAuth == nil {
			return "", fmt.Errorf("apikey_auth is required when auth_type=apikey")
		}
		b, err := json.Marshal(struct {
			Location string `json:"location"`
			Value    string `json:"value"`
		}{
			Location: req.APIKeyAuth.Location,
			Value:    req.APIKeyAuth.Value,
		})
		if err != nil {
			return "", err
		}
		return string(b), nil

	case source.AuthTypeSignature:
		if req.SignatureAuth == nil {
			return "", fmt.Errorf("signature_auth is required when auth_type=signature")
		}
		b, err := json.Marshal(struct {
			Secret    string `json:"secret"`
			Header    string `json:"header"`
			Algorithm string `json:"algorithm"`
			Encoding  string `json:"encoding"`
		}{
			Secret:    req.SignatureAuth.Secret,
			Header:    req.SignatureAuth.Header,
			Algorithm: req.SignatureAuth.Algorithm,
			Encoding:  req.SignatureAuth.Encoding,
		})
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("unsupported auth_type: %s", req.AuthType)
	}
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

func (s sourceService) List(ctx context.Context, page, limit int) ([]*source.Source, *response.Pagination, error) {
	//TODO implement me
	panic("implement me")
}

func (s sourceService) Update(ctx context.Context, id string, req source.UpdateRequest, currentUser *middleware.User) (*source.Source, error) {
	//TODO implement me
	panic("implement me")
}

func (s sourceService) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
