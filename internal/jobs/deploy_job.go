package jobs

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

func NewDeployTask(pod entities.Pod) *asynq.Task {
	payload, _ := json.Marshal(map[string]interface{}{
		"id": pod.Id,
		//"service": pod.Cluster,
		//"image":   pod.Image,
		//"port":    pod.Port,
	})

	return asynq.NewTask("deploy", payload)
}
