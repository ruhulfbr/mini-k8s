package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/moby/moby/client"
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/utils/fsUtils"
)

type DockerService struct {
	dockerConfig config.DockerConfig
}

func NewDockerService() *DockerService {
	return &DockerService{dockerConfig: config.GetDockerConfig()}
}

func (ds *DockerService) ValidateDockerContext(appName string, contextPath string) error {
	appPath := fsUtils.Join(ds.dockerConfig.ApplicationPath, appName)
	if !fsUtils.DirExists(appPath) {
		return appErrors.GirApplicationNotClonedYet
	}

	if !fsUtils.FileExists(fsUtils.Join(appPath, contextPath)) {
		return appErrors.DockerContextFileNotFound
	}

	return nil
}
func (ds *DockerService) BuildImage(service *entities.Service, appName string) (string, error) {
	imageTag := ds.generateImageTag(service.Name)
	buildContext := filepath.Join(ds.dockerConfig.ApplicationPath, appName, service.ContextPath)

	cmd := exec.CommandContext(
		context.Background(),
		"docker", "build",
		"-t", imageTag,
		buildContext,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("docker build failed for application: %s, Service: %s, Image Tag: %s,: %w", appName, service.Name, imageTag, err)
	}

	return imageTag, nil
}

func (ds *DockerService) generateImageTag(service string) string {
	return fmt.Sprintf(
		"%s-%s-%s:latest",
		ds.dockerConfig.ImageTagPrefix,
		service,
		uuid.NewV7,
	)
}

func (ds *DockerService) getContainerName(id, service string) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		ds.dockerConfig.ContainerNamePref,
		service,
		id,
	)
}

func (ds *DockerService) getContainerIP(ctx context.Context, cli *client.Client, containerID string) (string, error) {
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
