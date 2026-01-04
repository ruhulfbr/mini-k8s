package entities

import "time"

type ServiceType string
type ServiceStatus int

const (
	ServiceTypeHTTP       ServiceType   = "http"
	ServiceTypeWorker     ServiceType   = "worker"
	ServiceStatusActive   ServiceStatus = 1
	ServiceStatusInActive ServiceStatus = 0
)

type Service struct {
	Id            int64         `db:"id"              json:"id"`
	ApplicationId int64         `db:"application_id"  json:"application_id"`
	Name          string        `db:"name"            json:"name"`
	IP            string        `db:"ip"              json:"ip"`
	Port          int           `db:"port"            json:"port,omitempty"`
	ImageTag      string        `db:"image_tag"       json:"image_tag"`
	ContextPath   string        `db:"context_path"    json:"context_path"`
	Replicas      int           `db:"replicas"        json:"replicas"`
	Resources     string        `db:"resources"       json:"resources"`
	Path          string        `db:"path"            json:"path"`
	Type          ServiceType   `db:"type"            json:"type"`
	Status        ServiceStatus `db:"status"          json:"status"`
	LastBuildAt   *time.Time    `db:"last_build_at"   json:"last_build_at,omitempty"`
	CreatedAt     time.Time     `db:"created_at"      json:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at"      json:"updated_at"`
}
