package events

type ClusterEventAction string

const (
	ActionStarted ClusterEventAction = "started"
	ActionFailed  ClusterEventAction = "failed"
	ActionSuccess ClusterEventAction = "success"
)

type ClusterEvent string

const (
	EventBuildDockerImage    ClusterEvent = "build_image"
	EventPullDockerImage     ClusterEvent = "pull_image"
	EventDeploy              ClusterEvent = "deploy"
	EventRollingDeploy       ClusterEvent = "rolling_deploy"
	EventScaling             ClusterEvent = "scaling"
	EventStopClusterTraffic  ClusterEvent = "stop_traffic"
	EventStartClusterTraffic ClusterEvent = "start_traffic"
)

var ClusterEvents = map[ClusterEvent]map[ClusterEventAction]string{
	EventBuildDockerImage: {
		ActionStarted: "Docker image build started",
		ActionFailed:  "Docker image build failed",
		ActionSuccess: "Docker image build succeeded",
	},

	EventPullDockerImage: {
		ActionStarted: "Image pull started",
		ActionFailed:  "Image pull failed",
		ActionSuccess: "Image pull succeeded",
	},

	EventDeploy: {
		ActionStarted: "Application deployment started",
		ActionFailed:  "Application deployment failed",
		ActionSuccess: "Application deployed successfully",
	},

	EventScaling: {
		ActionStarted: "Scaling started",
		ActionFailed:  "Scaling failed",
		ActionSuccess: "Scaling succeeded",
	},

	EventRollingDeploy: {
		ActionStarted: "Rolling update started",
		ActionFailed:  "Rolling update failed",
		ActionSuccess: "Rolling update completed successfully",
	},

	EventStopClusterTraffic: {
		ActionStarted: "Stoping cluster traffic started",
		ActionFailed:  "Stoping cluster traffic failed",
		ActionSuccess: "Stoping cluster traffic succeeded",
	},

	EventStartClusterTraffic: {
		ActionStarted: "Starting cluster traffic operation started",
		ActionFailed:  "Starting cluster traffic operation failed",
		ActionSuccess: "Starting cluster traffic operation succeeded",
	},
}

func GetEventMessage(event ClusterEvent, action ClusterEventAction) string {
	if actions, ok := ClusterEvents[event]; ok {
		if msg, ok := actions[action]; ok {
			return msg
		}
	}
	return "Unknown event/action"
}
