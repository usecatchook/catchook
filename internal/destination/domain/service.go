package domain

import (
	"context"

	"github.com/theotruvelot/catchook/pkg/response"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*Destination, error)
	GetByID(ctx context.Context, id string) (*Destination, error)
	List(ctx context.Context, req ListDestinationsRequest) ([]*DestinationListItem, *response.Pagination, error)
	Update(ctx context.Context, id string, req UpdateRequest) (*Destination, error)
	Delete(ctx context.Context, id string) error
}
