package service

import (
	"context"
	"fmt"

	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/pkg/cache"
	"github.com/theotruvelot/catchook/pkg/crypto"
	"github.com/theotruvelot/catchook/pkg/logger"
)

type userService struct {
	userRepo user.Repository
	cache    cache.Cache
	logger   logger.Logger
}

func NewUserService(userRepo user.Repository, cache cache.Cache, logger logger.Logger) user.Service {
	return &userService{
		userRepo: userRepo,
		cache:    cache,
		logger:   logger,
	}
}

func (s *userService) Create(ctx context.Context, req user.CreateRequest) (*user.User, error) {
	s.logger.Info(ctx, "Creating new user", logger.String("email", req.Email))

	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		s.logger.Error(ctx, "Failed to check email existence", logger.Error(err))
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	if exists {
		s.logger.Warn(ctx, "Email already exists", logger.String("email", req.Email))
		return nil, user.ErrEmailAlreadyExists
	}

	hashedPassword, err := crypto.Hash(req.Password)
	if err != nil {
		s.logger.Error(ctx, "Failed to hash password", logger.Error(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &user.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		s.logger.Error(ctx, "Failed to create user", logger.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	cacheKey := cache.BuildKey(cache.KeyUserProfile, newUser.ID)
	s.cache.SetJSON(ctx, cacheKey, newUser, cache.TTLUserProfile)

	s.logger.Info(ctx, "User created successfully",
		logger.Int("user_id", newUser.ID),
		logger.String("email", newUser.Email),
	)

	newUser.Sanitize()
	return newUser, nil
}

func (s *userService) GetByID(ctx context.Context, id int) (*user.User, error) {
	s.logger.Debug(ctx, "Getting user by ID", logger.Int("user_id", id))

	cacheKey := cache.BuildKey(cache.KeyUserProfile, id)
	var cachedUser user.User
	if err := s.cache.GetJSON(ctx, cacheKey, &cachedUser); err == nil {
		s.logger.Debug(ctx, "User found in cache", logger.Int("user_id", id))
		cachedUser.Sanitize()
		return &cachedUser, nil
	}

	foundUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn(ctx, "User not found", logger.Int("user_id", id), logger.Error(err))
		return nil, user.ErrUserNotFound
	}

	s.cache.SetJSON(ctx, cacheKey, foundUser, cache.TTLUserProfile)

	s.logger.Debug(ctx, "User retrieved from database", logger.Int("user_id", id))

	foundUser.Sanitize()
	return foundUser, nil
}

func (s *userService) Update(ctx context.Context, id int, req user.UpdateRequest) (*user.User, error) {
	s.logger.Info(ctx, "Updating user", logger.Int("user_id", id))

	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn(ctx, "User not found for update", logger.Int("user_id", id))
		return nil, user.ErrUserNotFound
	}

	existingUser.FirstName = req.FirstName
	existingUser.LastName = req.LastName

	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		s.logger.Error(ctx, "Failed to update user", logger.Int("user_id", id), logger.Error(err))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	cacheKey := cache.BuildKey(cache.KeyUserProfile, id)
	s.cache.Delete(ctx, cacheKey)

	s.logger.Info(ctx, "User updated successfully", logger.Int("user_id", id))

	existingUser.Sanitize()
	return existingUser, nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	s.logger.Info(ctx, "Deleting user", logger.Int("user_id", id))

	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn(ctx, "User not found for deletion", logger.Int("user_id", id))
		return user.ErrUserNotFound
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, "Failed to delete user", logger.Int("user_id", id), logger.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	cacheKey := cache.BuildKey(cache.KeyUserProfile, id)
	s.cache.Delete(ctx, cacheKey)

	s.logger.Info(ctx, "User deleted successfully", logger.Int("user_id", id))
	return nil
}

func (s *userService) ChangePassword(ctx context.Context, id int, req user.ChangePasswordRequest) error {
	s.logger.Info(ctx, "Changing user password", logger.Int("user_id", id))

	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn(ctx, "User not found for password change", logger.Int("user_id", id))
		return user.ErrUserNotFound
	}

	valid, err := crypto.Verify(req.CurrentPassword, existingUser.Password)
	if err != nil || !valid {
		s.logger.Warn(ctx, "Invalid current password", logger.Int("user_id", id))
		return user.ErrInvalidPassword
	}

	hashedPassword, err := crypto.Hash(req.NewPassword)
	if err != nil {
		s.logger.Error(ctx, "Failed to hash new password", logger.Int("user_id", id), logger.Error(err))
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, id, hashedPassword); err != nil {
		s.logger.Error(ctx, "Failed to update password", logger.Int("user_id", id), logger.Error(err))
		return fmt.Errorf("failed to update password: %w", err)
	}

	cacheKey := cache.BuildKey(cache.KeyUserProfile, id)
	s.cache.Delete(ctx, cacheKey)

	s.logger.Info(ctx, "Password changed successfully", logger.Int("user_id", id))
	return nil
}
