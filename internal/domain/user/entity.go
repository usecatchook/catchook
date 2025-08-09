package user

import (
	"time"

	"github.com/theotruvelot/catchook/storage/postgres/generated"
)

type User struct {
	ID        string             `json:"id"`
	Email     string             `json:"email"`
	Role      generated.UserRole `json:"role"`
	Password  string             `json:"-"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	IsActive  bool               `json:"is_active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) Sanitize() {
	u.Password = ""
}
