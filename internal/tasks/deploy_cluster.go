package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type DeployClusterPayload struct {
	Cluster      entities.Cluster
	ClusterBuild entities.ClusterBuild
}

func DeployClusterTask(deployClusterPayload *DeployClusterPayload) *asynq.Task {
	payload, _ := json.Marshal(deployClusterPayload)

	return asynq.NewTask(DeployCluster, payload)
}
