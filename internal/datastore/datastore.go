package datastore

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/ruhulfbr/mini-k8s/internal/config"
)

type Datastore struct {
	DB *badger.DB
}

func NewDatastore(cfg *config.Config) *Datastore {
	opts := badger.DefaultOptions(cfg.Badger.DataSource).WithLogger(nil)
	DB, err := badger.Open(opts)
	if err != nil {

		log.Fatalf("failed to open BadgerDB: %v", err)
	}

	return &Datastore{DB: DB}
}

func (d *Datastore) Close() {
	d.DB.Close()
}
