package service

import (
	"context"
	"fmt"

	"github.com/theotruvelot/catchook/internal/platform/storage/postgres/generated"
	user "github.com/theotruvelot/catchook/internal/user/domain"
	"github.com/theotruvelot/catchook/pkg/cache"
	"github.com/theotruvelot/catchook/pkg/crypto"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/response"
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
		Role:      generated.UserRole(req.Role),
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		s.logger.Error(ctx, "Failed to create user", logger.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	newUser.Sanitize()
	return newUser, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*user.User, error) {
	s.logger.Debug(ctx, "Getting user by ID", logger.String("user_id", id))

	cacheKey := cache.BuildKey(cache.KeyUserProfile, id)
	var cachedUser user.User
	if err := s.cache.GetJSON(ctx, cacheKey, &cachedUser); err == nil {
		s.logger.Debug(ctx, "User found in cache", logger.String("user_id", id))
		cachedUser.Sanitize()
		return &cachedUser, nil
	}

	foundUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn(ctx, "User not found", logger.String("user_id", id), logger.Error(err))
		return nil, user.ErrUserNotFound
	}

	s.cache.SetJSON(ctx, cacheKey, foundUser, cache.TTLUserProfile)

	s.logger.Debug(ctx, "User retrieved from database", logger.String("user_id", id))

	foundUser.Sanitize()
	return foundUser, nil
}

func (s *userService) Update(ctx context.Context, id string, req user.UpdateRequest, currentUser *user.CurrentUser) (*user.User, error) {
	s.logger.Info(ctx, "Updating user", logger.String("user_id", id))

	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn(ctx, "User not found for update", logger.String("user_id", id))
		return nil, user.ErrUserNotFound
	}

	if currentUser == nil {
		return nil, user.ErrInsufficientPermissions
	}

	if currentUser.ID != id && currentUser.Role != "admin" {
		s.logger.Warn(ctx, "Non-admin trying to update another user",
			logger.String("current_user_id", currentUser.ID),
			logger.String("target_user_id", id))
		return nil, user.ErrInsufficientPermissions
	}

	existingUser.FirstName = req.FirstName
	existingUser.LastName = req.LastName

	if req.Role != "" {
		if currentUser.Role != "admin" {
			s.logger.Warn(ctx, "Non-admin trying to update role", logger.String("user_id", id))
			return nil, user.ErrInsufficientPermissions
		}
		existingUser.Role = generated.UserRole(req.Role)
	}

	if err := s.userRepo.Update(ctx, existingUser); err != nil {
		s.logger.Error(ctx, "Failed to update user", logger.String("user_id", id), logger.Error(err))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	cacheKey := cache.BuildKey(cache.KeyUserProfile, id)
	s.cache.Delete(ctx, cacheKey)

	s.logger.Info(ctx, "User updated successfully", logger.String("user_id", id))

	existingUser.Sanitize()
	return existingUser, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	s.logger.Info(ctx, "Deleting user", logger.String("user_id", id))

	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn(ctx, "User not found for deletion", logger.String("user_id", id))
		return user.ErrUserNotFound
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, "Failed to delete user", logger.String("user_id", id), logger.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	cacheKey := cache.BuildKey(cache.KeyUserProfile, id)
	s.cache.Delete(ctx, cacheKey)

	s.logger.Info(ctx, "User deleted successfully", logger.String("user_id", id))
	return nil
}

func (s *userService) ChangePassword(ctx context.Context, id string, req user.ChangePasswordRequest) error {
	s.logger.Info(ctx, "Changing user password", logger.String("user_id", id))

	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn(ctx, "User not found for password change", logger.String("user_id", id))
		return user.ErrUserNotFound
	}

	valid, err := crypto.Verify(req.CurrentPassword, existingUser.Password)
	if err != nil || !valid {
		s.logger.Warn(ctx, "Invalid current password", logger.String("user_id", id))
		return user.ErrInvalidPassword
	}

	hashedPassword, err := crypto.Hash(req.NewPassword)
	if err != nil {
		s.logger.Error(ctx, "Failed to hash new password", logger.String("user_id", id), logger.Error(err))
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, id, hashedPassword); err != nil {
		s.logger.Error(ctx, "Failed to update password", logger.String("user_id", id), logger.Error(err))
		return fmt.Errorf("failed to update password: %w", err)
	}

	cacheKey := cache.BuildKey(cache.KeyUserProfile, id)
	s.cache.Delete(ctx, cacheKey)

	s.logger.Info(ctx, "Password changed successfully", logger.String("user_id", id))
	return nil
}

func (s *userService) List(ctx context.Context, page, limit int) ([]*user.User, *response.Pagination, error) {
	s.logger.Debug(ctx, "Listing users", logger.Int("page", page), logger.Int("limit", limit))

	users, meta, err := s.userRepo.List(ctx, page, limit)
	if err != nil {
		s.logger.Error(ctx, "Failed to list users", logger.Error(err))
		return nil, nil, fmt.Errorf("failed to list users: %w", err)
	}

	for _, u := range users {
		u.Sanitize()
	}

	s.logger.Debug(ctx, "Users listed successfully",
		logger.Int("count", len(users)),
		logger.Int("page", page),
		logger.Int("limit", limit))
	return users, meta, nil
}
