package source

import (
	"context"

	"github.com/theotruvelot/catchook/pkg/response"
)

type Repository interface {
	Create(ctx context.Context, user *Source) error
	GetByID(ctx context.Context, id string) (*Source, error)
	List(ctx context.Context, page, limit int) ([]*Source, *response.Pagination, error)
	Update(ctx context.Context, user *Source) error
	Delete(ctx context.Context, id string) error
}
