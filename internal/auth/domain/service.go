package auth

import "context"

type Service interface {
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	RefreshSession(ctx context.Context, sessionID string) (*SessionResponse, error)
	Logout(ctx context.Context, sessionID string) error
}
