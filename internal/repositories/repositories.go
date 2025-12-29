package repositories

import (
	"github.com/dgraph-io/badger/v3"
)

type Repositories struct {
	PodRepository *PodRepository
}

func InitRepositories(DB *badger.DB) *Repositories {
	return &Repositories{
		PodRepository: NewPodRepository(DB),
	}
}
