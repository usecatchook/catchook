package source

type Source struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Protocol    string `json:"protocol"`
	AuthType    string `json:"auth_type"`
	AuthConfig  string `json:"auth_config"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
