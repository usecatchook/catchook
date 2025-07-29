package service

import (
	"context"
	"fmt"

	"github.com/theotruvelot/catchook/internal/domain/auth"
	"github.com/theotruvelot/catchook/internal/domain/user"
	"github.com/theotruvelot/catchook/pkg/crypto"
	"github.com/theotruvelot/catchook/pkg/jwt"
	"github.com/theotruvelot/catchook/pkg/logger"
)

type authService struct {
	userRepo    user.Repository
	userService user.Service
	jwtManager  jwt.Manager
	logger      logger.Logger
}

func NewAuthService(
	userRepo user.Repository,
	userService user.Service,
	jwtManager jwt.Manager,
	logger logger.Logger,
) auth.Service {
	return &authService{
		userRepo:    userRepo,
		userService: userService,
		jwtManager:  jwtManager,
		logger:      logger,
	}
}

func (s *authService) Login(ctx context.Context, req auth.LoginRequest) (*auth.AuthResponse, error) {
	s.logger.Info(ctx, "User login attempt", logger.String("email", req.Email))

	foundUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn(ctx, "Login failed - user not found", logger.String("email", req.Email))
		return nil, auth.ErrInvalidCredentials
	}

	if !foundUser.IsActive {
		s.logger.Warn(ctx, "Login failed - user inactive",
			logger.String("email", req.Email),
			logger.Int("user_id", foundUser.ID),
		)
		return nil, user.ErrUserInactive
	}

	valid, err := crypto.Verify(req.Password, foundUser.Password)
	if err != nil || !valid {
		s.logger.Warn(ctx, "Login failed - invalid password",
			logger.String("email", req.Email),
			logger.Int("user_id", foundUser.ID),
		)
		return nil, auth.ErrInvalidCredentials
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(foundUser.ID)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate access token",
			logger.Int("user_id", foundUser.ID),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(foundUser.ID)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate refresh token",
			logger.Int("user_id", foundUser.ID),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	s.logger.Info(ctx, "User logged in successfully",
		logger.Int("user_id", foundUser.ID),
		logger.String("email", foundUser.Email),
	)

	return &auth.AuthResponse{
		User: foundUser.ToResponse(),
		Tokens: &auth.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *authService) Register(ctx context.Context, req auth.RegisterRequest) (*auth.AuthResponse, error) {
	s.logger.Info(ctx, "User registration attempt", logger.String("email", req.Email))

	createReq := user.CreateRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	newUser, err := s.userService.Create(ctx, createReq)
	if err != nil {
		s.logger.Error(ctx, "Registration failed",
			logger.String("email", req.Email),
			logger.Error(err),
		)
		return nil, err
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(newUser.ID)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate access token for new user",
			logger.Int("user_id", newUser.ID),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(newUser.ID)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate refresh token for new user",
			logger.Int("user_id", newUser.ID),
			logger.Error(err),
		)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	s.logger.Info(ctx, "User registered successfully",
		logger.Int("user_id", newUser.ID),
		logger.String("email", newUser.Email),
	)

	return &auth.AuthResponse{
		User: newUser.ToResponse(),
		Tokens: &auth.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req auth.RefreshTokenRequest) (*auth.TokenPair, error) {
	s.logger.Debug(ctx, "Token refresh attempt")

	newAccessToken, err := s.jwtManager.RefreshAccessToken(ctx, req.RefreshToken)
	if err != nil {
		s.logger.Warn(ctx, "Token refresh failed", logger.Error(err))
		return nil, auth.ErrInvalidToken
	}

	claims, err := s.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		s.logger.Warn(ctx, "Invalid refresh token", logger.Error(err))
		return nil, auth.ErrInvalidToken
	}

	tokenPair := &auth.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: req.RefreshToken,
		TokenType:    "Bearer",
	}

	s.logger.Info(ctx, "Token refreshed successfully", logger.Int("user_id", claims.UserID))

	return tokenPair, nil
}
