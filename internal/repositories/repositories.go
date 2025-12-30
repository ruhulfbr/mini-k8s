package repositories

import (
	"database/sql"
)

type Repositories struct {
	ApplicationRepository  *ApplicationRepository
	ServiceRepository      *ServiceRepository
	BuildHistoryRepository *BuildHistoryRepository
	PodRepository          *PodRepository
}

func InitRepositories(DB *sql.DB) *Repositories {
	return &Repositories{
		ApplicationRepository:  NewApplicationRepository(DB),
		ServiceRepository:      NewServiceRepository(DB),
		BuildHistoryRepository: NewBuildHistoryRepository(DB),
		PodRepository:          NewPodRepository(DB),
	}
}
