package domain

import (
	"context"

	"github.com/theotruvelot/catchook/pkg/response"
)

type Repository interface {
	Create(ctx context.Context, destination *Destination) error
	GetByID(ctx context.Context, id string) (*Destination, error)
	GetByName(ctx context.Context, name string) (*Destination, error)
	List(ctx context.Context, req ListDestinationsRequest) ([]*DestinationListItem, *response.Pagination, error)
	Update(ctx context.Context, id, name, description string, destType DestinationType, config string, isActive bool, delaySeconds, retryAttempts int32) (*Destination, error)
	Delete(ctx context.Context, id string) error
}
