package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type BuildDockerImagePayload struct {
	ApplicationId   int64
	ApplicationName string
	ClusterId       int64
	ClusterName     string
	Version         string
	BuildConfig     entities.ClusterBuildConfig
}

func BuildDockerImageTask(BuildDockerImagePayload *BuildDockerImagePayload) *asynq.Task {
	payload, _ := json.Marshal(BuildDockerImagePayload)

	return asynq.NewTask(BuildDockerImage, payload)
}
