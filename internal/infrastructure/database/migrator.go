package database

import (
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(db *sqlx.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	applied, err := appliedMigrations(db)
	if err != nil {
		return err
	}

	files, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	ranAny := false
	for _, f := range files {
		version := strings.Split(f.Name(), "_")[0]

		if applied[version] {
			continue
		}

		log.Printf("Running migration: %s at %s", f.Name(), time.DateTime)

		sqlBytes, err := migrationFiles.ReadFile("migrations/" + f.Name())
		if err != nil {
			return err
		}

		if err := applyMigration(db, version, string(sqlBytes)); err != nil {
			return fmt.Errorf("migration %s failed: %w", f.Name(), err)
		}

		log.Printf("Migration applied: %s at %s", f.Name(), time.DateTime)
		ranAny = true
	}

	if !ranAny {
		log.Println("Nothing to migrate")
	}

	return nil
}

func ensureMigrationsTable(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func appliedMigrations(db *sqlx.DB) (map[string]bool, error) {
	rows, err := db.Query(`SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = true
	}
	return applied, nil
}

func applyMigration(db *sqlx.DB, version, sql string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(sql); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(
		`INSERT INTO schema_migrations (version) VALUES (?)`,
		version,
	); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
