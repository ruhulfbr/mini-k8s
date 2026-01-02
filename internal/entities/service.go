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
	Id            int64
	ApplicationId int64
	Name          string
	IP            string
	Port          *int
	ImageTag      string
	ContextPath   string
	Replicas      int
	Resources     string // JSON
	Path          string
	Type          ServiceType
	Status        ServiceStatus
	LastBuildAt   *time.Time
}
