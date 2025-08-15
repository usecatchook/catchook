package health

type StatusResponse struct {
	Status         string            `json:"status"`
	Version        string            `json:"version"`
	FirstTimeSetup bool              `json:"is_first_time_setup"`
	Message        string            `json:"message"`
	Services       map[string]string `json:"services"`
}
