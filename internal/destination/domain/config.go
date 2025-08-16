package domain

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type HTTPAuthType string

const (
	HTTPAuthTypeNone   HTTPAuthType = "none"
	HTTPAuthTypeBasic  HTTPAuthType = "basic"
	HTTPAuthTypeBearer HTTPAuthType = "bearer"
	HTTPAuthTypeAPIKey HTTPAuthType = "apikey"
)

type HTTPMethod string

const (
	HTTPMethodGET    HTTPMethod = "GET"
	HTTPMethodPOST   HTTPMethod = "POST"
	HTTPMethodPUT    HTTPMethod = "PUT"
	HTTPMethodPATCH  HTTPMethod = "PATCH"
	HTTPMethodDELETE HTTPMethod = "DELETE"
)

type HTTPContentType string

const (
	HTTPContentTypeJSON HTTPContentType = "application/json"
	HTTPContentTypeXML  HTTPContentType = "application/xml"
	HTTPContentTypeForm HTTPContentType = "application/x-www-form-urlencoded"
	HTTPContentTypeText HTTPContentType = "text/plain"
)

type HTTPAuth struct {
	Type     HTTPAuthType `json:"type" validate:"required,oneof=none basic bearer apikey"`
	Username string       `json:"username,omitempty"`
	Password string       `json:"password,omitempty"`
	Token    string       `json:"token,omitempty"`
	APIKey   string       `json:"api_key,omitempty"`
	Header   string       `json:"header,omitempty"`
}

type HTTPConfig struct {
	URL         string          `json:"url" validate:"required,url"`
	Method      HTTPMethod      `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	ContentType HTTPContentType `json:"content_type" validate:"omitempty,oneof=application/json application/xml application/x-www-form-urlencoded text/plain"`
	Timeout     int             `json:"timeout" validate:"omitempty,min=1,max=300"`
	Auth        *HTTPAuth       `json:"auth,omitempty"`
}

func (h *HTTPConfig) Validate() error {
	parsedURL, err := url.Parse(h.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	if h.Auth != nil {
		if err := h.Auth.Validate(); err != nil {
			return fmt.Errorf("auth validation failed: %w", err)
		}
	}

	if h.Timeout == 0 {
		h.Timeout = 30
	}

	if h.ContentType == "" {
		h.ContentType = HTTPContentTypeJSON
	}

	return nil
}

func (a *HTTPAuth) Validate() error {
	switch a.Type {
	case HTTPAuthTypeNone:
	case HTTPAuthTypeBasic:
		if a.Username == "" {
			return fmt.Errorf("username is required for basic auth")
		}
		if a.Password == "" {
			return fmt.Errorf("password is required for basic auth")
		}
	case HTTPAuthTypeBearer:
		if a.Token == "" {
			return fmt.Errorf("token is required for bearer auth")
		}
	case HTTPAuthTypeAPIKey:
		if a.APIKey == "" {
			return fmt.Errorf("api_key is required for apikey auth")
		}
		if a.Header == "" {
			return fmt.Errorf("header is required for apikey auth")
		}
		if strings.TrimSpace(a.Header) == "" {
			return fmt.Errorf("header name cannot be empty")
		}
	default:
		return fmt.Errorf("unsupported auth type: %s", a.Type)
	}
	return nil
}

func (h *HTTPConfig) ToMap() (map[string]interface{}, error) {
	data, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func HTTPConfigFromMap(data map[string]interface{}) (*HTTPConfig, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var config HTTPConfig
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func MergeHeaders(defaultHeaders, userHeaders map[string]string) map[string]string {
	merged := make(map[string]string)

	for k, v := range defaultHeaders {
		merged[k] = v
	}

	for k, v := range userHeaders {
		merged[k] = v
	}

	return merged
}
