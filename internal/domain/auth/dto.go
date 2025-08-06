package auth

import (
	"github.com/theotruvelot/catchook/internal/domain/user"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User    *user.UserResponse `json:"user"`
	Session *SessionResponse   `json:"session"`
}

type SessionResponse struct {
	SessionID string `json:"session_id"`
}
