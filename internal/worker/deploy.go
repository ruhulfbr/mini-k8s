package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

func (w *Worker) HandleDeploy() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		return nil
		//var payload entities.Pod
		//if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		//	return err
		//}
		//
		//fmt.Println("[Deploy] payload:", payload)
		//
		//ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
		//defer cancel()
		//
		//cli, err := client.New(client.FromEnv)
		//if err != nil {
		//	return fmt.Errorf("[Deploy] Docker client error: %v", err)
		//}
		//defer cli.Close()
		//
		//imageTag, err := w.buildImage(payload)
		//if err != nil {
		//	return fmt.Errorf("[Deploy] image build failed: %v", err)
		//}
		//
		//containerName := w.getContainerName(payload.Id, payload.Cluster)
		//
		//resp, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		//	Config: &container.Config{
		//		Image: imageTag,
		//	},
		//	HostConfig: &container.HostConfig{
		//		NetworkMode: "bridge",
		//	},
		//	NetworkingConfig: &network.NetworkingConfig{},
		//	Name:             containerName,
		//})
		//if err != nil {
		//	fmt.Println("Failed to create container on deploy worker:", err)
		//	return err
		//}
		//
		//// Start container
		//_, err = cli.ContainerStart(ctx, resp.Id, client.ContainerStartOptions{})
		//if err != nil {
		//	fmt.Println("Failed to start container on deploy worker:", err)
		//	return err
		//}
		//
		//ip, err := w.getContainerIP(cli, ctx, resp.Id)
		//if err != nil {
		//	return fmt.Errorf("[Deploy] get IP failed: %v", err)
		//}
		//
		//pod := entities.NewRunningPod(payload, ip)
		//
		//return w.podRepo.PutPod(ctx, pod)
	}
}
