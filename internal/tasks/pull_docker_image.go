package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type PullDockerImagePayload struct {
	ApplicationName string
	Version         string
	Cluster         entities.Cluster
}

func PullDockerImageTask(pullPayload *PullDockerImagePayload) *asynq.Task {
	payload, _ := json.Marshal(pullPayload)

	return asynq.NewTask(PullDockerImage, payload)
}
