package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	ApplicationRepository        *ApplicationRepository
	ClusterRepository            *ClusterRepository
	ClusterBuildConfigRepository *ClusterBuildConfigRepository
	ClusterBuildRepository       *ClusterBuildRepository
	PodRepository                *PodRepository
	ClusterEventRepository       *ClusterEventRepository
}

func InitRepositories(DB *sqlx.DB) *Repositories {
	return &Repositories{
		ApplicationRepository:        NewApplicationRepository(DB),
		ClusterRepository:            NewClusterRepository(DB),
		ClusterBuildConfigRepository: NewClusterBuildConfigRepository(DB),
		ClusterBuildRepository:       NewClusterBuildRepository(DB),
		PodRepository:                NewPodRepository(DB),
		ClusterEventRepository:       NewClusterEventRepository(DB),
	}
}
