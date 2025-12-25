package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/hibiken/asynq"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type DeployPayload struct {
	ID      string `json:"id"`
	Image   string `json:"image"`
	Service string `json:"service"`
	Port    int    `json:"port"`
}

type TerminatePayload struct {
	ID      string `json:"id"`
	Service string `json:"service"`
}

func StartWorker(ds *datastore.Datastore, redisAddr string) {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("deploy", handleDeploy(ds))
	mux.HandleFunc("terminate", handleTerminate(ds))

	if err := server.Run(mux); err != nil {
		log.Fatalf("worker failed: %v", err)
	}
}

// ---------------- DEPLOY ----------------

func handleDeploy(ds *datastore.Datastore) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload DeployPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		fmt.Println("Deploy payload on worker", payload)

		// timeout context for docker operations
		ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
		defer cancel()

		cli, err := client.New(client.FromEnv)
		if err != nil {
			fmt.Println("Failed to create client on deploy worker:", err)
			return err
		}
		defer cli.Close()

		//// Pull image
		//reader, err := cli.ImagePull(ctx, payload.Image, client.ImagePullOptions{})
		//if err != nil {
		//	fmt.Println("Image Pull error on deploy worker:", err)
		//	return err
		//}
		//io.Copy(io.Discard, reader)
		//reader.Close()

		imageTag, err := buildImageFromDockerFile(ctx, payload)
		if err != nil {
			fmt.Println("Failed to build image on deploy worker:", err)
			return err
		}

		containerName := generateContainerName(payload.ID, payload.Service)

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

		// Inspect to get IP
		inspect, err := cli.ContainerInspect(ctx, resp.ID, client.ContainerInspectOptions{})
		if err != nil {
			fmt.Println("Failed to inspect container on deploy worker:", err)
			return err
		}

		ip := ""
		for _, net := range inspect.Container.NetworkSettings.Networks {
			ip = net.IPAddress.String()
			break
		}

		if ip == "" {
			return fmt.Errorf("container has no IP address")
		}

		// Update pod status in BadgerDB
		pod := entities.Pod{
			ID:      payload.ID,
			Service: payload.Service,
			Image:   payload.Image,
			Port:    payload.Port,
			Status:  entities.PodRunning,
			IP:      ip,
		}

		return ds.PutPod(ctx, pod)
	}
}

// ---------------- TERMINATE ----------------

func handleTerminate(ds *datastore.Datastore) asynq.HandlerFunc {
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

		containerName := generateContainerName(payload.ID, payload.Service)

		_, err = cli.ContainerRemove(ctx, containerName, client.ContainerRemoveOptions{Force: true, RemoveVolumes: true})

		if err != nil {
			log.Printf("failed to remove container %s: %v", containerName, err)

			return err
		}

		return ds.DeletePod(ctx, payload.Service, payload.ID)
	}
}

func buildImageFromDockerFile(ctx context.Context, payload DeployPayload) (string, error) {
	imageTag := generateImageTag(payload.Service)

	buildContext := "./nodes/" + payload.Service

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"build",
		"-t", imageTag,
		buildContext,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return imageTag, cmd.Run()
}

func generateImageTag(serviceName string) string {
	return fmt.Sprintf(
		"%s-%s:latest",
		os.Getenv("IMAGE_TAG_PREFIX"),
		serviceName,
	)
}

func generateContainerName(id string, service string) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		os.Getenv("CONTAINER_NAME_PREFIX"),
		service,
		id,
	)
}
