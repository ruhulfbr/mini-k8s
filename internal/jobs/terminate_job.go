package jobs

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

func NewTerminateTask(pod entities.Pod) *asynq.Task {
	payload, _ := json.Marshal(map[string]string{
		"id": "s",
		// "service": pod.Cluster,
	})

	return asynq.NewTask("terminate", payload)
}
