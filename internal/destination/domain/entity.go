package domain

import (
	"time"
)

type Destination struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	DestinationType DestinationType `json:"destination_type"`
	Config          string          `json:"config"`
	IsActive        bool            `json:"is_active"`
	DelaySeconds    int32           `json:"delay_seconds"`
	RetryAttempts   int32           `json:"retry_attempts"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
