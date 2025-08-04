package user

import (
	"context"

	"github.com/theotruvelot/catchook/internal/middleware"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	Update(ctx context.Context, id int, req UpdateRequest, currentUser *middleware.User) (*User, error)
	Delete(ctx context.Context, id int) error
	ChangePassword(ctx context.Context, id int, req ChangePasswordRequest) error
}
