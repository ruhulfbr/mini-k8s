package worker

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

func (w *Worker) buildImage(payload entities.Pod) (string, error) {
	return "", nil

	//imageTag := w.generateImageTag(payload.Service)
	//buildContext := "./nodes/" + payload.Service
	//
	//cmd := exec.CommandContext(
	//	context.Background(),
	//	"docker", "build",
	//	"-t", imageTag,
	//	buildContext,
	//)
	//
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//
	//return imageTag, cmd.Run()
}

func (w *Worker) generateImageTag(service string) string {
	return fmt.Sprintf("%s-%s:latest", w.cfg.Docker.ImageTagPrefix, service)
}

func (w *Worker) getContainerName(id, service string) string {
	return fmt.Sprintf("%s-%s-%s", w.cfg.Docker.ContainerNamePref, service, id)
}

func (w *Worker) getContainerIP(cli *client.Client, ctx context.Context, containerID string) (string, error) {
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
