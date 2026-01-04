package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	ApplicationRepository  *ApplicationRepository
	ServiceRepository      *ServiceRepository
	BuildHistoryRepository *BuildHistoryRepository
	PodRepository          *PodRepository
}

func InitRepositories(DB *sqlx.DB) *Repositories {
	return &Repositories{
		ApplicationRepository:  NewApplicationRepository(DB),
		ServiceRepository:      NewServiceRepository(DB),
		BuildHistoryRepository: NewBuildHistoryRepository(DB),
		PodRepository:          NewPodRepository(DB),
	}
}
