package entities

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
