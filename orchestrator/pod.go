package orchestrator

import (
	"fmt"
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

/*
PodKey design:

	pods/<service>/<podID>
	This allows:
	- List pods by service
	- Efficient deletes
*/
func PodKey(service, id string) string {
	return fmt.Sprintf("pods/%s/%s", service, id)
}
