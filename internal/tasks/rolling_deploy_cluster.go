package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type RollingDeployClusterPayload struct {
	Cluster      entities.Cluster
	ClusterBuild entities.ClusterBuild
}

func RollingDeployClusterTask(deployPayload *RollingDeployClusterPayload) *asynq.Task {
	payload, _ := json.Marshal(deployPayload)

	return asynq.NewTask(RollingDeployCluster, payload)
}
