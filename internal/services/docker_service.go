package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/moby/moby/client"
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/utils"
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

func (ds *DockerService) BuildImage(buildConfig *entities.ClusterBuildConfig, appName string, clusterName string) (string, error) {
	dockerContextPath := ""
	if buildConfig.DockerContextPath != "." {
		dockerContextPath = buildConfig.DockerContextPath
	}

	buildContext := filepath.Join(
		ds.dockerConfig.ClusterPath,
		ds.clusterDir(appName, clusterName),
		dockerContextPath,
	)
	dockerfilePath := filepath.Join(buildContext, buildConfig.DockerfileName)

	if !fsUtils.FileExists(dockerfilePath) {
		return "", appErrors.DockerFileNotFound
	}

	ctx := context.Background()
	imageTag := ds.generateImageTag(appName, clusterName)

	cmd := exec.CommandContext(
		ctx,
		"docker", "build",
		"-f", dockerfilePath,
		"-t", imageTag,
		buildContext,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	err := cmd.Run()
	if err != nil {
		logger.Error(ctx, "Docker build failed", err,
			"application", appName,
			"Cluster", clusterName,
			"Image Tag", imageTag,
			"Error Details", stderr.String(),
		)

		return "", appErrors.DockerFailedToBuildImage
	}

	return imageTag, nil
}

func (ds *DockerService) PullImageWithTag(appName string, cluster *entities.Cluster) (string, error) {
	ctx := context.Background()

	imageTag := *cluster.CurrentImageTag

	reader, err := ds.cli.ImagePull(ctx, imageTag, client.ImagePullOptions{})
	if err != nil {
		logger.Error(ctx, "Failed to pull docker image", err,
			"imageTag", imageTag,
			"New Image Tag", imageTag,
			"Application", appName,
			"Cluster", cluster.Name,
		)
		return "", appErrors.DockerFailedToPullImage
	}
	defer reader.Close()

	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		logger.Error(ctx, "Error reading pull response", err,
			"imageTag", imageTag,
			"Application", appName,
			"Cluster", cluster.Name,
		)
		return "", appErrors.DockerFailedReadPullResponse
	}

	// Tag the image locally
	newImageTag := imageTag + "-" + utils.UniqueId()
	_, err = ds.cli.ImageTag(ctx, client.ImageTagOptions{
		Source: imageTag,
		Target: newImageTag,
	})
	if err != nil {
		logger.Error(ctx, "Failed to tag docker image", err,
			"imageTag", imageTag,
			"New Image Tag", newImageTag,
			"Application", appName,
			"Cluster", cluster.Name,
		)
		return "", appErrors.DockerFailedAddNewTagToImage
	}

	return newImageTag, nil
}

func (ds *DockerService) generateImageTag(appName string, clusterName string) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		appName,
		clusterName,
		utils.UniqueId(),
	)
}

func (ds *DockerService) getContainerName(id, appName string, clusterName string) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		appName,
		clusterName,
		id,
	)
}

func (ds *DockerService) getContainerIP(ctx context.Context, cli *client.Client, containerId string) (string, error) {
	inspect, err := cli.ContainerInspect(ctx, containerId, client.ContainerInspectOptions{})
	if err != nil {
		return "", err
	}

	ip := ""
	for _, net := range inspect.Container.NetworkSettings.Networks {
		ip = net.IPAddress.String()
		break
	}

	if ip == "" {
		logger.Error(ctx, "Failed to get container IP", err, "containerId", containerId)
		return "", appErrors.DockerContainerHasNoIP
	}

	return ip, nil
}

func (ds *DockerService) clusterDir(appName string, clusterName string) string {
	return appName + "-" + clusterName
}
