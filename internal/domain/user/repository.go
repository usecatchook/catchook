package user

import (
	"context"

	"github.com/theotruvelot/catchook/pkg/response"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, page, limit int) ([]*User, *response.Pagination, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	EmailExists(ctx context.Context, email string) (bool, error)
	UpdatePassword(ctx context.Context, userID string, hashedPassword string) error
	CountUsers(ctx context.Context) (int64, error)
}
