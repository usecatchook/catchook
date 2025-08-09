package service

import (
	"context"
	"fmt"

	"github.com/theotruvelot/catchook/internal/domain/auth"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/pkg/crypto"
	"github.com/theotruvelot/catchook/pkg/logger"
	"github.com/theotruvelot/catchook/pkg/session"
)

type authService struct {
	userRepo       user.Repository
	sessionManager session.Manager
	logger         logger.Logger
}

func NewAuthService(
	userRepo user.Repository,
	sessionManager session.Manager,
	logger logger.Logger,
) auth.Service {
	return &authService{
		userRepo:       userRepo,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

func (s *authService) Login(ctx context.Context, req auth.LoginRequest) (*auth.AuthResponse, error) {
	s.logger.Info(ctx, "User login attempt", logger.String("email", req.Email))

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn(ctx, "Login failed - user not found", logger.String("email", req.Email))
		return nil, auth.ErrInvalidCredentials
	}

	if !user.IsActive {
		s.logger.Warn(ctx, "Login failed - user inactive",
			logger.String("email", req.Email),
			logger.String("user_id", user.ID),
		)
		return nil, auth.ErrUserInactive
	}

	valid, err := crypto.Verify(req.Password, user.Password)
	if err != nil || !valid {
		s.logger.Warn(ctx, "Login failed - invalid password",
			logger.String("email", req.Email),
			logger.String("user_id", user.ID),
		)
		return nil, auth.ErrInvalidCredentials
	}

	sessionID, err := s.sessionManager.CreateSession(ctx, user.ID, string(user.Role))
	if err != nil {
		s.logger.Error(ctx, "Failed to create session",
			logger.String("user_id", user.ID),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info(ctx, "User logged in successfully",
		logger.String("user_id", user.ID),
		logger.String("email", user.Email),
	)

	return &auth.AuthResponse{
		User: user.ToResponse(),
		Session: &auth.SessionResponse{
			SessionID: sessionID,
		},
	}, nil
}

func (s *authService) RefreshSession(ctx context.Context, sessionID string) (*auth.SessionResponse, error) {
	s.logger.Debug(ctx, "Session refresh attempt")

	sess, err := s.sessionManager.ValidateSession(ctx, sessionID)
	if err != nil {
		s.logger.Warn(ctx, "Invalid session", logger.Error(err))
		return nil, auth.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, sess.UserID)
	if err != nil {
		s.logger.Warn(ctx, "User not found during session refresh",
			logger.String("user_id", sess.UserID),
			logger.Error(err))
		return nil, auth.ErrInvalidToken
	}

	if !user.IsActive {
		s.logger.Warn(ctx, "Session refresh failed - user inactive",
			logger.String("user_id", user.ID),
		)
		return nil, auth.ErrUserInactive
	}

	err = s.sessionManager.RefreshSession(ctx, sessionID)
	if err != nil {
		s.logger.Error(ctx, "Failed to refresh session",
			logger.String("user_id", user.ID),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	s.logger.Info(ctx, "Session refreshed successfully",
		logger.String("user_id", user.ID),
		logger.String("role", string(user.Role)))

	return &auth.SessionResponse{
		SessionID: sessionID,
	}, nil
}

func (s *authService) Logout(ctx context.Context, sessionID string) error {
	s.logger.Debug(ctx, "Processing logout")

	err := s.sessionManager.DeleteSession(ctx, sessionID)
	if err != nil {
		s.logger.Error(ctx, "Failed to delete session during logout",
			logger.String("session_id", sessionID),
			logger.Error(err),
		)
		return fmt.Errorf("failed to delete session: %w", err)
	}

	s.logger.Info(ctx, "User logged out successfully",
		logger.String("session_id", sessionID))

	return nil
}
