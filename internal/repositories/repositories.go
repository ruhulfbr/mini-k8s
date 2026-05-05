package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	ContextRepository            *ContextRepository
	ClusterRepository            *ClusterRepository
	ClusterBuildConfigRepository *ClusterBuildConfigRepository
	ClusterBuildRepository       *ClusterBuildRepository
	PodRepository                *PodRepository
	ClusterEventRepository       *ClusterEventRepository
}

func InitRepositories(DB *sqlx.DB) *Repositories {
	return &Repositories{
		ContextRepository:            NewContextRepository(DB),
		ClusterRepository:            NewClusterRepository(DB),
		ClusterBuildConfigRepository: NewClusterBuildConfigRepository(DB),
		ClusterBuildRepository:       NewClusterBuildRepository(DB),
		PodRepository:                NewPodRepository(DB),
		ClusterEventRepository:       NewClusterEventRepository(DB),
	}
}
