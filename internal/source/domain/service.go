package source

import (
	"context"

	"github.com/theotruvelot/catchook/pkg/response"
)

// CurrentUser represents the authenticated user context required by the domain layer
// Kept minimal to avoid coupling with transport/middleware packages
type CurrentUser struct {
	ID   string
	Role string
}

type Service interface {
	Create(ctx context.Context, req CreateRequest, currentUser *CurrentUser) (*Source, error)
	GetByID(ctx context.Context, id string) (*Source, error)
	List(ctx context.Context, page, limit int) ([]*SourceResponse, *response.Pagination, error)
	Update(ctx context.Context, id string, req UpdateRequest, currentUser *CurrentUser) (*Source, error)
	Delete(ctx context.Context, id string) error
}
