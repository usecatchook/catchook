package source

import (
	"encoding/json"
	"fmt"
	"time"
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
	Name          string               `json:"name" validate:"required,min=2,max=50"`
	Description   string               `json:"description" validate:"omitempty,max=255"`
	Protocol      string               `json:"protocol" validate:"required,oneof=http grpc mqtt websocket"`
	AuthType      AuthType             `json:"auth_type" validate:"required,oneof=none basic bearer apikey signature"`
	BasicAuth     *BasicAuthConfig     `json:"basic_auth" validate:"required_if=AuthType basic"`
	BearerAuth    *BearerAuthConfig    `json:"bearer_auth" validate:"required_if=AuthType bearer"`
	APIKeyAuth    *APIKeyAuthConfig    `json:"apikey_auth" validate:"required_if=AuthType apikey"`
	SignatureAuth *SignatureAuthConfig `json:"signature_auth" validate:"required_if=AuthType signature"`
}

type BasicAuthConfig struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

type BearerAuthConfig struct {
	Token string `json:"token" validate:"required,min=1"`
}

type APIKeyAuthConfig struct {
	Location string `json:"location" validate:"required,min=1"`
	Value    string `json:"value" validate:"required,min=1"`
}

type SignatureAuthConfig struct {
	Secret    string `json:"secret" validate:"required,min=1"`
	Header    string `json:"header" validate:"required,min=1"`
	Algorithm string `json:"algorithm" validate:"required,oneof=sha-1 sha-256 sha-512 md5"`
	Encoding  string `json:"encoding" validate:"required,oneof=base64 base64url hex"`
}

type UpdateRequest struct {
	Name          string               `json:"name" validate:"required,min=2,max=50"`
	Description   string               `json:"description" validate:"omitempty,max=255"`
	Protocol      string               `json:"protocol" validate:"required,oneof=http grpc mqtt websocket"`
	BasicAuth     *BasicAuthConfig     `json:"basic_auth" validate:"required_if=AuthType basic"`
	BearerAuth    *BearerAuthConfig    `json:"bearer_auth" validate:"required_if=AuthType bearer"`
	APIKeyAuth    *APIKeyAuthConfig    `json:"apikey_auth" validate:"required_if=AuthType apikey"`
	SignatureAuth *SignatureAuthConfig `json:"signature_auth" validate:"required_if=AuthType signature"`
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

type ListUsersRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	Limit    int    `query:"limit" validate:"omitempty,min=1"`
	Protocol string `query:"protocol" validate:"omitempty,oneof=http grpc mqtt websocket"`
	Search   string `query:"search" validate:"omitempty,min=2,max=50"`
	OrderBy  string `query:"order_by" validate:"omitempty,oneof=name created_at updated_at"`
	Order    string `query:"order" validate:"omitempty,oneof=asc desc"`
	IsActive bool   `query:"is_active" validate:"omitempty,boolean"`
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
