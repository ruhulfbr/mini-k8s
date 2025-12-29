package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

func (w *Worker) HandleDeploy() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload entities.Pod
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		fmt.Println("[Deploy] payload:", payload)

		ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
		defer cancel()

		cli, err := client.New(client.FromEnv)
		if err != nil {
			return fmt.Errorf("[Deploy] Docker client error: %v", err)
		}
		defer cli.Close()

		imageTag, err := w.buildImage(payload)
		if err != nil {
			return fmt.Errorf("[Deploy] image build failed: %v", err)
		}

		containerName := w.getContainerName(payload.ID, payload.Service)

		resp, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
			Config: &container.Config{
				Image: imageTag,
			},
			HostConfig: &container.HostConfig{
				NetworkMode: "bridge",
			},
			NetworkingConfig: &network.NetworkingConfig{},
			Name:             containerName,
		})
		if err != nil {
			fmt.Println("Failed to create container on deploy worker:", err)
			return err
		}

		// Start container
		_, err = cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
		if err != nil {
			fmt.Println("Failed to start container on deploy worker:", err)
			return err
		}

		ip, err := w.getContainerIP(cli, ctx, resp.ID)
		if err != nil {
			return fmt.Errorf("[Deploy] get IP failed: %v", err)
		}

		pod := entities.NewRunningPod(payload, ip)

		return w.podRepo.PutPod(ctx, pod)
	}
}
