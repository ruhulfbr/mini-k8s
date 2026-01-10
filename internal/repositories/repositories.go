package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	ApplicationRepository        *ApplicationRepository
	ServiceRepository            *ServiceRepository
	ServiceBuildConfigRepository *ServiceBuildConfigRepository
	BuildHistoryRepository       *BuildHistoryRepository
	PodRepository                *PodRepository
}

func InitRepositories(DB *sqlx.DB) *Repositories {
	return &Repositories{
		ApplicationRepository:        NewApplicationRepository(DB),
		ServiceRepository:            NewServiceRepository(DB),
		ServiceBuildConfigRepository: NewServiceBuildConfigRepository(DB),
		BuildHistoryRepository:       NewBuildHistoryRepository(DB),
		PodRepository:                NewPodRepository(DB),
	}
}
