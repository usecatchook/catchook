package source

import (
	"time"
)

type CreateRequest struct {
	Name          string               `json:"name" validate:"required,min=2,max=50"`
	Description   string               `json:"description" validate:"omitempty,max=255"`
	Protocol      string               `json:"protocol" validate:"required,oneof=http grpc mqtt websocket"`
	AuthType      string               `json:"auth_type" validate:"required,oneof=none basic bearer apikey signature"`
	BasicAuth     *BasicAuthConfig     `json:"basic_auth,omitempty" validate:"omitempty,required_if=AuthType basic,dive"`
	BearerAuth    *BearerAuthConfig    `json:"bearer_auth,omitempty" validate:"omitempty,required_if=AuthType bearer,dive"`
	APIKeyAuth    *APIKeyAuthConfig    `json:"apikey_auth,omitempty" validate:"omitempty,required_if=AuthType apikey,dive"`
	SignatureAuth *SignatureAuthConfig `json:"signature_auth,omitempty" validate:"omitempty,required_if=AuthType signature,dive"`
}

type BasicAuthConfig struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

type BearerAuthConfig struct {
	Token string `json:"token" validate:"required,min=1"`
}

type APIKeyAuthConfig struct {
	Key      string `json:"key" validate:"required,min=1"`
	Location string `json:"location" validate:"required,oneof=header query"`
	Name     string `json:"name" validate:"required,min=1"`
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
	BasicAuth     *BasicAuthConfig     `json:"basic_auth,omitempty" validate:"omitempty,required_if=AuthType basic,dive"`
	BearerAuth    *BearerAuthConfig    `json:"bearer_auth,omitempty" validate:"omitempty,required_if=AuthType bearer,dive"`
	APIKeyAuth    *APIKeyAuthConfig    `json:"apikey_auth,omitempty" validate:"omitempty,required_if=AuthType apikey,dive"`
	SignatureAuth *SignatureAuthConfig `json:"signature_auth,omitempty" validate:"omitempty,required_if=AuthType signature,dive"`
}

type SourceResponse struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	Protocol      string               `json:"protocol"`
	AuthType      string               `json:"auth_type"`
	IsActive      bool                 `json:"is_active"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	BasicAuth     *BasicAuthConfig     `json:"basic_auth,omitempty"`
	BearerAuth    *BearerAuthConfig    `json:"bearer_auth,omitempty"`
	APIKeyAuth    *APIKeyAuthConfig    `json:"apikey_auth,omitempty"`
	SignatureAuth *SignatureAuthConfig `json:"signature_auth,omitempty"`
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
