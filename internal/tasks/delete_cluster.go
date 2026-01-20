package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

type DeleteClusterPayload struct {
	ClusterId int64
}

func DeleteClusterTask(deleteClusterPayload *DeleteClusterPayload) *asynq.Task {
	payload, _ := json.Marshal(deleteClusterPayload)

	return asynq.NewTask(DeleteCluster, payload)
}
