package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/theotruvelot/catchook/pkg/response"
)

type DestinationType string

const (
	DestinationTypeHTTP     DestinationType = "http"
	DestinationTypeRabbitMQ DestinationType = "rabbitmq"
	DestinationTypeDatabase DestinationType = "database"
	DestinationTypeFile     DestinationType = "file"
	DestinationTypeQueue    DestinationType = "queue"
	DestinationTypeCLI      DestinationType = "cli"
)

type CreateRequest struct {
	Name            string                 `json:"name" validate:"required,min=2,max=100"`
	Description     string                 `json:"description" validate:"omitempty,max=255"`
	DestinationType DestinationType        `json:"destination_type" validate:"required,oneof=http rabbitmq database file queue cli"`
	Config          map[string]interface{} `json:"config" validate:"omitempty"`
	DelaySeconds    int32                  `json:"delay_seconds" validate:"omitempty,min=0"`
	RetryAttempts   int32                  `json:"retry_attempts" validate:"omitempty,min=0"`
}

type UpdateRequest struct {
	Name            string                 `json:"name" validate:"omitempty,min=2,max=100"`
	Description     string                 `json:"description" validate:"omitempty,max=255"`
	DestinationType DestinationType        `json:"destination_type" validate:"omitempty,oneof=http rabbitmq database file queue cli"`
	Config          map[string]interface{} `json:"config" validate:"omitempty"`
	IsActive        bool                   `json:"is_active" validate:"omitempty"`
	DelaySeconds    int32                  `json:"delay_seconds" validate:"omitempty,min=0"`
	RetryAttempts   int32                  `json:"retry_attempts" validate:"omitempty,min=0"`
}

type DestinationResponse struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	DestinationType string                 `json:"destination_type"`
	Config          map[string]interface{} `json:"config,omitempty"`
	IsActive        bool                   `json:"is_active"`
	DelaySeconds    int32                  `json:"delay_seconds"`
	RetryAttempts   int32                  `json:"retry_attempts"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type ListDestinationsRequest struct {
	Page            int    `query:"page" validate:"omitempty,min=1"`
	Limit           int    `query:"limit" validate:"omitempty,min=1"`
	DestinationType string `query:"destination_type" validate:"omitempty,oneof=http rabbitmq database file queue cli"`
	Search          string `query:"search" validate:"omitempty,min=2,max=50"`
	OrderBy         string `query:"order_by" validate:"omitempty,oneof=name created_at updated_at is_active"`
	Order           string `query:"order" validate:"omitempty,oneof=asc desc"`
	IsActive        *bool  `query:"is_active" validate:"omitempty"`
}

type DestinationListItem struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DestinationType string    `json:"destination_type"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ListDestinationsResponse struct {
	Destinations []*DestinationListItem `json:"data"`
	Pagination   *response.Pagination   `json:"pagination"`
}

func (d *Destination) ToResponse() (*DestinationResponse, error) {
	resp := &DestinationResponse{
		ID:              d.ID,
		Name:            d.Name,
		Description:     d.Description,
		DestinationType: string(d.DestinationType),
		IsActive:        d.IsActive,
		DelaySeconds:    d.DelaySeconds,
		RetryAttempts:   d.RetryAttempts,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}

	if d.Config != "" {
		var cfg map[string]interface{}
		if err := json.Unmarshal([]byte(d.Config), &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal config: %w", err)
		}
		resp.Config = cfg
	}

	return resp, nil
}

func (d *Destination) ToListItem() *DestinationListItem {
	return &DestinationListItem{
		Name:            d.Name,
		Description:     d.Description,
		DestinationType: string(d.DestinationType),
		IsActive:        d.IsActive,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

func ToResponses(list []*Destination) ([]*DestinationResponse, error) {
	resp := make([]*DestinationResponse, 0, len(list))
	for _, item := range list {
		r, err := item.ToResponse()
		if err != nil {
			return nil, err
		}
		resp = append(resp, r)
	}
	return resp, nil
}

func ToListItems(list []*Destination) []*DestinationListItem {
	resp := make([]*DestinationListItem, 0, len(list))
	for _, item := range list {
		resp = append(resp, item.ToListItem())
	}
	return resp
}
