package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/moby/moby/client"
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/utils/fsUtils"
)

type DockerService struct {
	dockerConfig config.DockerConfig
	cli          *client.Client
}

func NewDockerService() *DockerService {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		panic(fmt.Sprintf("failed to create docker client: %v", err))
	}

	return &DockerService{
		dockerConfig: config.GetDockerConfig(),
		cli:          cli,
	}
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

func (ds *DockerService) BuildImage(buildConfig *entities.ClusterBuildConfig, appName string, clusterName string) (string, error) {
	imageTag := ds.generateImageTag(clusterName)

	dockerContextPath := ""
	if buildConfig.DockerContextPath != "." {
		dockerContextPath = buildConfig.DockerContextPath
	}

	buildContext := filepath.Join(ds.dockerConfig.ApplicationPath, appName, dockerContextPath)
	dockerfilePath := filepath.Join(buildContext, buildConfig.DockerfileName)

	cmd := exec.CommandContext(
		context.Background(),
		"docker", "build",
		"-f", dockerfilePath,
		"-t", imageTag,
		buildContext,
	)

	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(
			"docker build failed (app=%s cluster=%s tag=%s): %w",
			appName,
			clusterName,
			imageTag,
			err,
		)
	}

	return imageTag, nil
}

func (ds *DockerService) PullImage(imageTag string) error {
	ctx := context.Background()

	reader, err := ds.cli.ImagePull(ctx, imageTag, client.ImagePullOptions{})
	if err != nil {
		logger.Error(ctx, "Failed to pull docker image", err)
		return err
	}
	defer reader.Close()

	_, err = io.Copy(io.Discard, reader)
	return err
}

func (ds *DockerService) generateImageTag(clusterName string) string {
	uuId, _ := uuid.NewV7()

	return fmt.Sprintf(
		"%s-%s-%s",
		ds.dockerConfig.ImageTagPrefix,
		clusterName,
		uuId.String(),
	)
}

func (ds *DockerService) getContainerName(id, clusterName string) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		ds.dockerConfig.ContainerNamePref,
		clusterName,
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
