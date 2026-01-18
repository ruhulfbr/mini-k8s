package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ClusterBuildConfigRepository struct {
	db *sqlx.DB
}

func NewClusterBuildConfigRepository(db *sqlx.DB) *ClusterBuildConfigRepository {
	return &ClusterBuildConfigRepository{db: db}
}

func (r *ClusterBuildConfigRepository) Create(cfg *entities.ClusterBuildConfig) error {
	query := `
	INSERT INTO cluster_build_configs (
		cluster_id, git_repo, git_branch,
		docker_context_path, dockerfile_name
	)
	VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(
		query,
		cfg.ClusterId,
		cfg.GitRepo,
		cfg.GitBranch,
		cfg.DockerContextPath,
		cfg.DockerfileName,
	)
	return err
}

func (r *ClusterBuildConfigRepository) Update(cfg *entities.ClusterBuildConfig) error {
	query := `
	UPDATE cluster_build_configs
	SET
		git_repo = ?,
		git_branch = ?,
		docker_context_path = ?,
		dockerfile_name = ?,
		updated_at = CURRENT_TIMESTAMP
	WHERE cluster_id = ?
	`

	_, err := r.db.Exec(
		query,
		cfg.GitRepo,
		cfg.GitBranch,
		cfg.DockerContextPath,
		cfg.DockerfileName,
		cfg.ClusterId,
	)

	return err
}

func (r *ClusterBuildConfigRepository) GetByClusterId(clusterId int64) (*entities.ClusterBuildConfig, error) {
	var cfg entities.ClusterBuildConfig

	query := `SELECT * FROM cluster_build_configs WHERE cluster_id = ? LIMIT 1`

	err := r.db.Get(&cfg, query, clusterId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cfg, nil
}

func (r *ClusterBuildConfigRepository) Delete(clusterId int64) error {
	_, err := r.db.Exec(`DELETE FROM cluster_build_configs WHERE cluster_id = ?`, clusterId)

	if err != nil {
		return err
	}

	return nil
}
