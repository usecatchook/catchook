package user

import (
	"github.com/theotruvelot/catchook/storage/postgres/generated"
	"time"
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
