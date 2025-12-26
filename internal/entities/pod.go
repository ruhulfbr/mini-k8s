package entities

import (
	"github.com/google/uuid"
	"github.com/ruhulfbr/mini-k8s/internal/http/requests"
)

type PodStatus string

const (
	PodPending PodStatus = "Pending"
	PodRunning PodStatus = "Running"
)

type Pod struct {
	ID      string    `json:"id"`
	Service string    `json:"service"`
	Image   string    `json:"image"`
	Port    int       `json:"port"`
	Status  PodStatus `json:"status"`
	IP      string    `json:"ip"`
}

func NewPendingPod(req requests.ScaleRequest) Pod {
	return Pod{
		ID:      uuid.NewString(),
		Service: req.ServiceName,
		Image:   req.Image,
		Port:    req.Port,
		Status:  PodPending,
	}
}

func NewRunningPod(payload Pod, ip string) Pod {
	return Pod{
		ID:      payload.ID,
		Service: payload.Service,
		Image:   payload.Image,
		Port:    payload.Port,
		Status:  PodRunning,
		IP:      ip,
	}
}
