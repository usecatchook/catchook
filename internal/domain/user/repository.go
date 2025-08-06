package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
	EmailExists(ctx context.Context, email string) (bool, error)
	UpdatePassword(ctx context.Context, userID int, hashedPassword string) error
	CountUsers(ctx context.Context) (int64, error)
}
