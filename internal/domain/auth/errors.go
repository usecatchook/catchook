package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserInactive       = errors.New("user is inactive")
)
