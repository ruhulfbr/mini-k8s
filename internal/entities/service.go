package entities

import "time"

type ServiceType string
type ServiceStatus int
type DeployMode int

const (
	ServiceTypeHTTP       ServiceType   = "http"
	ServiceTypeWorker     ServiceType   = "worker"
	ServiceStatusActive   ServiceStatus = 1
	ServiceStatusInActive ServiceStatus = 0
	DeployModeImage       DeployMode    = 1
	DeployModeBuild       DeployMode    = 2
)

type Service struct {
	Id              int64         `db:"id"              json:"id"`
	ApplicationId   int64         `db:"application_id"  json:"application_id"`
	Name            string        `db:"name"            json:"name"`
	IP              string        `db:"ip"              json:"ip"`
	Port            int           `db:"port"              json:"port"`
	Replicas        int           `db:"replicas"        json:"replicas"`
	CPU             int           `db:"cpu"       json:"cpu"`
	Memory          int           `db:"memory"       json:"memory"`
	Path            string        `db:"path"            json:"path"`
	Type            ServiceType   `db:"type"            json:"type"`
	DeployMode      DeployMode    `db:"deploy_mode"          json:"deploy_mode"`
	Image           *string       `db:"image"             json:"image"`
	ContextPath     *string       `db:"context_path"    json:"context_path"`
	DockerFile      *string       `db:"docker_file"    json:"docker_file"`
	CurrentImageTag *string       `db:"current_image_tag" json:"current_image_tag"`
	CurrentVersion  *string       `db:"current_version"   json:"current_version"`
	Status          ServiceStatus `db:"status"          json:"status"`
	LastDeployedAt  *time.Time    `db:"last_deployed_at" json:"last_deployed_at"`
	CreatedAt       time.Time     `db:"created_at"      json:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at"      json:"updated_at"`
}
