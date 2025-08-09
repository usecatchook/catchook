package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theotruvelot/catchook/internal/domain/source"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/pkg/tracer"
	"github.com/theotruvelot/catchook/storage/postgres/generated"
)

type sourceRepository struct {
	db        *pgxpool.Pool
	queries   *generated.Queries
	appLogger logger.Logger
}

func NewSourceRepository(db *pgxpool.Pool, appLogger logger.Logger) source.Repository {
	return &sourceRepository{
		db:        db,
		queries:   generated.New(db),
		appLogger: appLogger,
	}
}

func (s sourceRepository) Create(ctx context.Context, source *source.Source) error {
	ctx, span := tracer.StartSpan(ctx, "source.repository.create")
	defer span.End()

	userId, err := uuid.Parse(source.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	result, err := s.queries.CreateSource(ctx,
		source.Name,
		userId,
		source.Description,
		generated.ProtocolType(source.Protocol),
		generated.AuthType(source.AuthType),
		[]byte(source.AuthConfig),
		source.IsActive,
	)

	if err != nil {
		s.appLogger.Error(ctx, "Failed to create source",
			logger.String("name", source.Name),
			logger.String("user_id", source.UserID),
			logger.Error(err),
		)
		span.RecordError(err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	source.ID = result.ID.String()
	source.CreatedAt = result.CreatedAt.Time
	source.UpdatedAt = result.UpdatedAt.Time

	return nil
}

func (s sourceRepository) GetByID(ctx context.Context, id string) (*source.Source, error) {
	ctx, span := tracer.StartSpan(ctx, "source.repository.get_by_id")
	defer span.End()

	uid, err := uuid.Parse(id)
	if err != nil {
		s.appLogger.Error(ctx, "Invalid UUID format",
			logger.String("source_id", id),
			logger.Error(err),
		)
		return nil, fmt.Errorf("invalid source ID format: %w", err)
	}
	result, err := s.queries.GetSourceByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		s.appLogger.Error(ctx, "Failed to get source by ID", logger.Error(err))
		return nil, fmt.Errorf("failed to get source by ID: %w", err)
	}

	return &source.Source{
		ID:          result.ID.String(),
		UserID:      result.UserID.String(),
		Name:        result.Name,
		Description: result.Description,
		Protocol:    string(result.Protocol),
		AuthType:    source.AuthType(result.AuthType),
		AuthConfig:  string(result.AuthConfig),
		IsActive:    false,
		CreatedAt:   result.CreatedAt.Time,
		UpdatedAt:   result.UpdatedAt.Time,
	}, nil
}

func (s sourceRepository) List(ctx context.Context, page, limit int) ([]*source.Source, *response.Pagination, error) {
	//TODO implement me
	panic("implement me")
}

func (s sourceRepository) Update(ctx context.Context, user *source.Source) error {
	//TODO implement me
	panic("implement me")
}

func (s sourceRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s sourceRepository) GetByName(ctx context.Context, name string) (*source.Source, error) {
	ctx, span := tracer.StartSpan(ctx, "source.repository.get_by_name")
	defer span.End()

	result, err := s.queries.GetSourceByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &source.Source{
		ID:          result.ID.String(),
		Name:        result.Name,
		UserID:      result.UserID.String(),
		Description: result.Description,
		Protocol:    string(result.Protocol),
		AuthType:    source.AuthType(result.AuthType),
		AuthConfig:  string(result.AuthConfig),
		IsActive:    false,
		CreatedAt:   result.CreatedAt.Time,
		UpdatedAt:   result.UpdatedAt.Time,
	}, nil
}
