package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ScaleClusterPayload struct {
	Cluster      entities.Cluster
	ClusterBuild entities.ClusterBuild
	Replicas     int
}

func ScaleClusterTask(scaleClusterPayload *ScaleClusterPayload) *asynq.Task {
	payload, _ := json.Marshal(scaleClusterPayload)

	return asynq.NewTask(ScaleCluster, payload)
}
