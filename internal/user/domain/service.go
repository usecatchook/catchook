package user

import (
	"context"

	"github.com/theotruvelot/catchook/pkg/response"
)

// CurrentUser represents the authenticated user executing the action
// Minimal to avoid coupling with presentation layer
type CurrentUser struct {
	ID   string
	Role string
}

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context, page, limit int) ([]*User, *response.Pagination, error)
	Update(ctx context.Context, id string, req UpdateRequest, currentUser *CurrentUser) (*User, error)
	Delete(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, id string, req ChangePasswordRequest) error
}
