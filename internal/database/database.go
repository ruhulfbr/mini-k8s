package database

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
	"github.com/ruhulfbr/mini-k8s/internal/config"
)

type Database struct {
	DB *sql.DB
}

func NewDatastore(cfg *config.Config) *Database {
	db, err := sql.Open("sqlite", cfg.SQLite.DataSource)
	if err != nil {
		log.Fatalf("failed to open BadgerDB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping sqlite db: %v", err)
	}

	log.Println("Connected to sqlite database")

	// Strong defaults for control-plane workloads
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA busy_timeout = 5000;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			log.Fatalf("failed to apply pragma: %v", err)
		}
	}

	if err := RunMigrations(db); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	log.Println("Migrations applied successfully")

	return &Database{DB: db}
}

func (d *Database) Close() {
	d.DB.Close()
}
