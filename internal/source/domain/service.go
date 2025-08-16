package source

import (
	"context"

	"github.com/theotruvelot/catchook/pkg/response"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*Source, error)
	GetByID(ctx context.Context, id string) (*Source, error)
	List(ctx context.Context, page, limit int) ([]*SourceResponse, *response.Pagination, error)
	Update(ctx context.Context, id string, req UpdateRequest) (*Source, error)
	Delete(ctx context.Context, id string) error
}
