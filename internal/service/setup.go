package service

import (
	"context"
	"fmt"

	"github.com/theotruvelot/catchook/internal/domain/setup"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/pkg/crypto"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/storage/postgres/generated"
)

type setupService struct {
	userRepo user.Repository
	logger   logger.Logger
}

func NewSetupService(userRepo user.Repository, logger logger.Logger) setup.Service {
	return &setupService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s setupService) CreateAdminUser(ctx context.Context, req setup.CreateAdminUserRequest) error {
	s.logger.Info(ctx, "Creating first admin user", logger.String("email", req.Email))

	count, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to count users", logger.Error(err))
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count > 0 {
		s.logger.Warn(ctx, "Setup already completed, users exist")
		return setup.ErrAdminUserAlreadyExists
	}

	hashedPassword, err := crypto.Hash(req.Password)
	if err != nil {
		s.logger.Error(ctx, "Failed to hash password", logger.Error(err))
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &user.User{
		Email:     req.Email,
		Role:      generated.UserRoleAdmin,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		s.logger.Error(ctx, "Failed to create user", logger.Error(err))
		return fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info(ctx, "Admin user created successfully", logger.String("user_id", newUser.ID), logger.String("email", newUser.Email))
	return nil
}
