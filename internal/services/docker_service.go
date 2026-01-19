package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
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

type ContainerInfo struct {
	IP            string
	ContainerID   string
	ContainerName string
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

	imageTag := *cluster.Image

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

func (ds *DockerService) DeployImage(cluster *entities.Cluster, buildInfo *entities.ClusterBuild) (*ContainerInfo, error) {
	ctx := context.Background()

	containerName := ds.generateContainerName(cluster.Name)
	resp, err := ds.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: buildInfo.ImageTag,
		},
		HostConfig: &container.HostConfig{
			NetworkMode: "bridge",
		},
		NetworkingConfig: &network.NetworkingConfig{},
		Name:             containerName,
	})
	if err != nil {
		logger.Error(ctx, "Failed to create container", err,
			"containerName", containerName,
			"Cluster", cluster.Name,
			"imageTag", buildInfo.ImageTag,
			"Version", buildInfo.Version,
		)
		return nil, appErrors.DockerFailedCreateContainer
	}

	// Start container
	_, err = ds.cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		logger.Error(ctx, "Failed to run container", err,
			"containerName", containerName,
			"Cluster", cluster.Name,
			"imageTag", buildInfo.ImageTag,
			"Version", buildInfo.Version,
		)
		return nil, appErrors.DockerFailedRunContainer
	}

	ip, err := ds.getContainerIP(ctx, resp.ID)
	if err != nil {
		logger.Error(ctx, "Failed to retrieve container IP", err,
			"containerId", resp.ID,
			"containerName", containerName,
			"Cluster", cluster.Name,
			"imageTag", buildInfo.ImageTag,
			"Version", buildInfo.Version,
		)
		return nil, err
	}

	return &ContainerInfo{
		IP:            ip,
		ContainerID:   resp.ID,
		ContainerName: containerName,
	}, nil
}

func (ds *DockerService) DeleteContainer(pod entities.Pod) error {
	ctx := context.Background()

	if _, err := ds.cli.ContainerRemove(ctx, pod.ContainerName, client.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}); err != nil {

		logger.Error(ctx, "Failed to remove container", err,
			"containerName", pod.ContainerName,
			"containerId", pod.ContainerName,
			"clusterId", pod.ClusterId,
		)

		return appErrors.DockerFailedDeleteContainer
	}

	return nil
}

func (ds *DockerService) generateImageTag(appName string, clusterName string) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		appName,
		clusterName,
		utils.UniqueId(),
	)
}

func (ds *DockerService) generateContainerName(clusterName string) string {
	return fmt.Sprintf(
		"%s-%s",
		clusterName,
		utils.UniqueId(),
	)
}

func (ds *DockerService) getContainerIP(ctx context.Context, containerId string) (string, error) {
	inspect, err := ds.cli.ContainerInspect(ctx, containerId, client.ContainerInspectOptions{})
	if err != nil {
		return "", err
	}

	ip := ""
	for _, net := range inspect.Container.NetworkSettings.Networks {
		ip = net.IPAddress.String()
		break
	}

	if ip == "" {
		return "", appErrors.DockerFailedRetrieveContainerIP
	}

	return ip, nil
}

func (ds *DockerService) clusterDir(appName string, clusterName string) string {
	return appName + "-" + clusterName
}
