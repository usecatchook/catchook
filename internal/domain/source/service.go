package source

import (
	"context"

	"github.com/theotruvelot/catchook/internal/middleware"
	"github.com/theotruvelot/catchook/pkg/response"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest, currentUser *middleware.User) (*Source, error)
	GetByID(ctx context.Context, id string) (*Source, error)
	List(ctx context.Context, page, limit int) ([]*Source, *response.Pagination, error)
	Update(ctx context.Context, id string, req UpdateRequest, currentUser *middleware.User) (*Source, error)
	Delete(ctx context.Context, id string) error
}
