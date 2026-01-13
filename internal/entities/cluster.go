package entities

import "time"

type ClusterType string
type ClusterStatus int
type DeployMode int

const (
	ClusterTypeHTTP       ClusterType   = "http"
	ClusterTypeWorker     ClusterType   = "worker"
	ClusterStatusActive   ClusterStatus = 1
	ClusterStatusInActive ClusterStatus = 0
	DeployModeImage       DeployMode    = 1
	DeployModeBuild       DeployMode    = 2
)

type Cluster struct {
	Id              int64             `db:"id"              json:"id"`
	ApplicationId   int64             `db:"application_id"  json:"application_id"`
	Name            string            `db:"name"            json:"name"`
	IP              string            `db:"ip"              json:"ip"`
	Port            int               `db:"port"              json:"port"`
	Replicas        int               `db:"replicas"        json:"replicas"`
	CPU             int               `db:"cpu"       json:"cpu"`
	Memory          int               `db:"memory"       json:"memory"`
	Path            string            `db:"path"            json:"path"`
	Type            ClusterType       `db:"type"            json:"type"`
	DeployMode      DeployMode        `db:"deploy_mode"          json:"deploy_mode"`
	Image           *string           `db:"image" json:"image"`
	Envs            map[string]string `db:"envs" json:"envs"`
	CurrentImageTag *string           `db:"current_image_tag" json:"current_image_tag"`
	CurrentVersion  *string           `db:"current_version"   json:"current_version"`
	Status          ClusterStatus     `db:"status"          json:"status"`
	LastDeployedAt  *time.Time        `db:"last_deployed_at" json:"last_deployed_at"`
	CreatedAt       time.Time         `db:"created_at"      json:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at"      json:"updated_at"`
}
