package workerServices

import (
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/services"
)

type Services struct {
	ClusterService *ClusterService
	GitService     *services.GitService
	DockerService  *services.DockerService
}

func InitWorkerServices(DB *sqlx.DB, asynqClient *asynq.Client) *Services {
	repos := repositories.InitRepositories(DB)

	gitService := services.NewGitService()
	dockerService := services.NewDockerService()

	return &Services{
		ClusterService: NewClusterService(
			repos.ClusterRepository,
			repos.ClusterBuildRepository,
			repos.PodRepository,
			gitService, dockerService,
		),
		GitService:    gitService,
		DockerService: dockerService,
	}
}
