package worker

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/moby/moby/client"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

func buildImage(payload entities.Pod) (string, error) {
	imageTag := generateImageTag(payload.Service)
	buildContext := "./nodes/" + payload.Service

	cmd := exec.CommandContext(
		context.Background(),
		"docker", "build",
		"-t", imageTag,
		buildContext,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return imageTag, cmd.Run()
}

func generateImageTag(service string) string {
	return fmt.Sprintf("%s-%s:latest", os.Getenv("IMAGE_TAG_PREFIX"), service)
}

func getContainerName(id, service string) string {
	return fmt.Sprintf("%s-%s-%s", os.Getenv("CONTAINER_NAME_PREFIX"), service, id)
}

func getContainerIP(cli *client.Client, ctx context.Context, containerID string) (string, error) {
	inspect, err := cli.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
	if err != nil {
		return "", err
	}

	ip := ""
	for _, net := range inspect.Container.NetworkSettings.Networks {
		ip = net.IPAddress.String()
		break
	}

	if ip == "" {
		return "", fmt.Errorf("container has no IP address")
	}

	return ip, nil
}
