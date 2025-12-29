package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/moby/moby/client"
)

type TerminatePayload struct {
	ID      string `json:"id"`
	Service string `json:"service"`
}

func (w *Worker) HandleTerminate() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload TerminatePayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		cli, err := client.New(client.FromEnv)
		if err != nil {
			return err
		}
		defer cli.Close()

		containerName := w.getContainerName(payload.ID, payload.Service)

		if _, err := cli.ContainerRemove(ctx, containerName, client.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		}); err != nil {
			log.Printf("[Terminate] failed to remove container %s: %v", containerName, err)
			return err
		}

		return w.podRepo.DeletePod(ctx, payload.Service, payload.ID)
	}
}
