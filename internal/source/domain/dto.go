package source

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/theotruvelot/catchook/pkg/response"
)

type AuthType string

const (
	AuthTypeNone      AuthType = "none"
	AuthTypeBasic     AuthType = "basic"
	AuthTypeBearer    AuthType = "bearer"
	AuthTypeApikey    AuthType = "apikey"
	AuthTypeSignature AuthType = "signature"
)

type CreateRequest struct {
	Name        string         `json:"name" validate:"required,min=2,max=50"`
	Description string         `json:"description" validate:"omitempty,max=255"`
	Protocol    string         `json:"protocol" validate:"required,oneof=http grpc mqtt websocket"`
	AuthType    AuthType       `json:"auth_type" validate:"required,oneof=none basic bearer apikey signature"`
	AuthConfig  map[string]any `json:"auth_config" validate:"omitempty"`
}

type UpdateRequest struct {
	Name        string         `json:"name" validate:"omitempty,min=2,max=50"`
	Description string         `json:"description" validate:"omitempty,max=255"`
	Protocol    string         `json:"protocol" validate:"omitempty,oneof=http grpc mqtt websocket"`
	AuthType    AuthType       `json:"auth_type" validate:"omitempty,oneof=none basic bearer apikey signature"`
	AuthConfig  map[string]any `json:"auth_config" validate:"omitempty"`
}

type SourceResponse struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Protocol   string         `json:"protocol"`
	AuthType   string         `json:"auth_type"`
	IsActive   bool           `json:"is_active"`
	AuthConfig map[string]any `json:"auth_config,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type ListSourcesRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	Limit    int    `query:"limit" validate:"omitempty,min=1"`
	Protocol string `query:"protocol" validate:"omitempty,oneof=http grpc mqtt websocket"`
	Search   string `query:"search" validate:"omitempty,min=2,max=50"`
	OrderBy  string `query:"order_by" validate:"omitempty,oneof=name created_at updated_at"`
	Order    string `query:"order" validate:"omitempty,oneof=asc desc"`
	IsActive bool   `query:"is_active" validate:"omitempty,boolean"`
}

type ListSourcesResponse struct {
	Sources    []*SourceResponse    `json:"data"`
	Pagination *response.Pagination `json:"pagination"`
}

func (s *Source) ToResponse() (*SourceResponse, error) {
	resp := &SourceResponse{
		ID:        s.ID,
		Name:      s.Name,
		Protocol:  s.Protocol,
		AuthType:  string(s.AuthType),
		IsActive:  s.IsActive,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}

	if s.AuthType == AuthTypeNone || s.AuthConfig == "" {
		return resp, nil
	}

	var cfg map[string]any
	if err := json.Unmarshal([]byte(s.AuthConfig), &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal auth config: %w", err)
	}
	resp.AuthConfig = cfg

	return resp, nil
}

func ToResponses(list []*Source) ([]*SourceResponse, error) {
	resp := make([]*SourceResponse, 0, len(list))
	for _, item := range list {
		r, err := item.ToResponse()
		if err != nil {
			return nil, err
		}
		resp = append(resp, r)
	}
	return resp, nil
}
