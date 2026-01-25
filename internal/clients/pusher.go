package clients

import (
	"fmt"

	"github.com/pusher/pusher-http-go/v5"
	"github.com/ruhulfbr/mini-k8s/internal/config"
)

var (
	ChannelClusterEvents = "channel_cluster_events"
	ClusterEvents        = "cluster_events"
)

type PusherClient struct {
	client *pusher.Client
}

func NewPusherClient() *PusherClient {
	pusherConfig := config.GetPusherConfig()

	return &PusherClient{
		client: &pusher.Client{
			AppID:   pusherConfig.AppId,
			Key:     pusherConfig.AppKey,
			Secret:  pusherConfig.AppSecret,
			Cluster: pusherConfig.AppCluster,
			Secure:  true,
		},
	}
}

func (pc *PusherClient) TriggerClusterEvent(data map[string]string) {
	err := pc.client.Trigger(ChannelClusterEvents, ClusterEvents, data)
	if err != nil {
		fmt.Println("err.Error()", err.Error())
	}
}
