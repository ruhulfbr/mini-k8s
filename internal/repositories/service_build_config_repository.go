package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ServiceBuildConfigRepository struct {
	db *sqlx.DB
}

func NewServiceBuildConfigRepository(db *sqlx.DB) *ServiceBuildConfigRepository {
	return &ServiceBuildConfigRepository{db: db}
}

func (r *ServiceBuildConfigRepository) Create(cfg *entities.ServiceBuildConfig) error {
	query := `
	INSERT INTO service_build_configs (
		service_id, git_repo, git_branch,
		docker_context_path, dockerfile_name
	)
	VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(
		query,
		cfg.ServiceId,
		cfg.GitRepo,
		cfg.GitBranch,
		cfg.DockerContextPath,
		cfg.DockerfileName,
	)
	return err
}

func (r *ServiceBuildConfigRepository) Update(cfg *entities.ServiceBuildConfig) error {
	query := `
	UPDATE service_build_configs
	SET
		git_repo = ?,
		git_branch = ?,
		docker_context_path = ?,
		dockerfile_name = ?,
		updated_at = CURRENT_TIMESTAMP
	WHERE service_id = ?
	`

	_, err := r.db.Exec(
		query,
		cfg.GitRepo,
		cfg.GitBranch,
		cfg.DockerContextPath,
		cfg.DockerfileName,
		cfg.ServiceId,
	)

	return err
}

func (r *ServiceBuildConfigRepository) GetByServiceId(serviceId int64) (*entities.ServiceBuildConfig, error) {
	var cfg entities.ServiceBuildConfig

	query := `SELECT * FROM service_build_configs WHERE service_id = ? LIMIT 1`

	err := r.db.Get(&cfg, query, serviceId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cfg, nil
}

func (r *ServiceBuildConfigRepository) Delete(serviceId int) error {
	res, err := r.db.Exec(`DELETE FROM service_build_configs WHERE service_id = ?`, serviceId)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
