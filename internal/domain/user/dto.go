package user

import (
	"time"

	"github.com/theotruvelot/catchook/pkg/response"
	"github.com/theotruvelot/catchook/storage/postgres/generated"
)

type CreateRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,password"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
}

type UpdateRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	Role      string `json:"role" validate:"omitempty,oneof=admin developer viewer"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,password"`
}

type UserResponse struct {
	ID        int                `json:"id"`
	Email     string             `json:"email"`
	Role      generated.UserRole `json:"role"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	FullName  string             `json:"full_name"`
	IsActive  bool               `json:"is_active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}

type ListUsersRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	Limit    int    `query:"limit" validate:"omitempty,min=1"`
	Role     string `query:"role" validate:"omitempty,oneof=admin developer viewer"`
	IsActive bool   `query:"is_active" validate:"omitempty,boolean"`
	OrderBy  string `query:"order_by" validate:"omitempty,oneof=first_name last_name role created_at updated_at"`
	Order    string `query:"order" validate:"omitempty,oneof=asc desc"`
	Search   string `query:"search" validate:"omitempty,min=2,max=50"`
}

type ListUsersResponse struct {
	Users      []*UserResponse      `json:"data"`
	Pagination *response.Pagination `json:"pagination"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		FullName:  u.FullName(),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
