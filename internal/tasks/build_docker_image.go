package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type BuildDockerImagePayload struct {
	ApplicationId   int64                       `json:"applicationId"`
	ApplicationName string                      `json:"applicationName"`
	ClusterId       int64                       `json:"clusterId"`
	ClusterName     string                      `json:"clusterName"`
	Version         string                      `json:"version"`
	BuildConfig     entities.ClusterBuildConfig `json:"buildConfig"`
}

func NewBuildDockerTaskTask(BuildDockerImagePayload *BuildDockerImagePayload) *asynq.Task {
	payload, _ := json.Marshal(BuildDockerImagePayload)

	return asynq.NewTask(BuildDockerImage, payload)
}
