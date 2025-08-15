package source

import (
	"time"
)

type Source struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Protocol    string    `json:"protocol"`
	AuthType    AuthType  `json:"auth_type"`
	AuthConfig  string    `json:"auth_config"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
