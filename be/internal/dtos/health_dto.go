package dtos

import "time"

// HealthStatus represents the overall health status response.
type HealthStatus struct {
	Status    string               `json:"status"`
	Timestamp time.Time            `json:"timestamp"`
	Database  *HealthComponent     `json:"database"`
	Redis     *HealthComponent     `json:"redis,omitempty"`
}

// HealthComponent represents the status of a single component.
type HealthComponent struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"`
}
