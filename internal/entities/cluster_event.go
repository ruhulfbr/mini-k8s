package entities

import "time"

type ClusterEventType string

const (
	EventClusterCreated         ClusterEventType = "cluster_created"
	EventClusterUpdated         ClusterEventType = "cluster_updated"
	EventClusterDeleted         ClusterEventType = "cluster_deleted"
	EventBuildStarted           ClusterEventType = "build_started"
	EventBuildSucceeded         ClusterEventType = "build_succeeded"
	EventBuildFailed            ClusterEventType = "build_failed"
	EventCodeCloning            ClusterEventType = "code_cloning"
	EventCodeCloned             ClusterEventType = "code_cloned"
	EventCodeCloningFailed      ClusterEventType = "code_cloning_failed"
	EventCodePulling            ClusterEventType = "code_pulling"
	EventCodePulled             ClusterEventType = "code_pulled"
	EventCodePullFailed         ClusterEventType = "code_pull_failed"
	EventPullImageStarted       ClusterEventType = "pull_image_started"
	EventPullImageFailed        ClusterEventType = "pull_image_failed"
	EventPulledImage            ClusterEventType = "pulled_image"
	ClusterDeployStarted        ClusterEventType = "cluster_deploy_started"
	ClusterDeployFailed         ClusterEventType = "cluster_deploy_failed"
	ClusterDeploySucceeded      ClusterEventType = "cluster_deploy_succeeded"
	EventPodStarted             ClusterEventType = "pod_started"
	EventPodReady               ClusterEventType = "pod_ready"
	EventPodRemoved             ClusterEventType = "pod_removed"
	EventPodFailed              ClusterEventType = "pod_failed"
	EventPodHealthCheckFailed   ClusterEventType = "pod_healthcheck_failed"
	EventScaleUp                ClusterEventType = "scale_up"
	EventScaleDown              ClusterEventType = "scale_down"
	EventRollingUpdateStart     ClusterEventType = "rolling_update_start"
	EventRollingUpdateFailed    ClusterEventType = "rolling_update_failed"
	EventRollingUpdateDone      ClusterEventType = "rolling_update_done"
	EventLoadBalancerStarted    ClusterEventType = "load-balancer_started"
	EventLoadBalancerStopped    ClusterEventType = "load-balancer_stopped"
	EventLoadBalancerPodAdded   ClusterEventType = "load-balancer_pod_added"
	EventLoadBalancerPodRemoved ClusterEventType = "load-balancer_pod_removed"
	EventClusterPaused          ClusterEventType = "cluster_paused"
	EventClusterResumed         ClusterEventType = "cluster_resumed"
)

type ClusterEvent struct {
	ID        int64            `db:"id" json:"id"`
	ClusterId int64            `db:"cluster_id" json:"cluster_id"`
	PodId     *int64           `db:"pod_id" json:"pod_id,omitempty"`
	Type      ClusterEventType `db:"type" json:"type"`
	Metadata  string           `db:"metadata" json:"metadata"`
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
}
