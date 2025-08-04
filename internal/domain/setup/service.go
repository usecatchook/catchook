package setup

import (
	"context"
)

type Service interface {
	CreateAdminUser(ctx context.Context, req CreateAdminUserRequest) error
}
