package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ClusterBuildRepository struct {
	db *sqlx.DB
}

type ClusterBuildRepositoryInterface interface {
	Create(history *entities.ClusterBuild) error
	GetByCluster(clusterId int64) ([]entities.ClusterBuild, error)
	GetRollbackImage(clusterId int64) (*string, error)
	Delete(id int64) error
}

func NewClusterBuildRepository(db *sqlx.DB) *ClusterBuildRepository {
	return &ClusterBuildRepository{db: db}
}

func (r *ClusterBuildRepository) GetByCluster(clusterId int64) ([]entities.ClusterBuild, error) {
	query := `SELECT * FROM cluster_builds WHERE cluster_id = ? ORDER BY id DESC`

	var clusterBuilds []entities.ClusterBuild
	if err := r.db.Select(&clusterBuilds, query, clusterId); err != nil {
		return nil, err
	}

	return clusterBuilds, nil
}

func (r *ClusterBuildRepository) Create(b *entities.ClusterBuild) error {
	return r.db.QueryRow(`
		INSERT INTO cluster_builds (cluster_id, version, image_tag)
		VALUES (?, ?, ?)
		RETURNING id, created_at`,
		b.ClusterId,
		b.Version,
		b.ImageTag,
	).Scan(&b.Id, &b.CreatedAt)
}

func (r *ClusterBuildRepository) GetById(clusterId int64, id int64) (*entities.ClusterBuild, error) {
	var bh entities.ClusterBuild
	err := r.db.Get(&bh, `SELECT * FROM cluster_builds WHERE cluster_id = ? and id = ?`, clusterId, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &bh, err
}

func (r *ClusterBuildRepository) ExistsByVersion(clusterId int64, version string) bool {
	var s entities.Cluster
	err := r.db.Get(&s, `SELECT id FROM cluster_builds WHERE cluster_id = ? and version = ?`, clusterId, version)

	if err != nil {
		return false
	}

	return true
}

func (r *ClusterBuildRepository) ExistsByImage(clusterId int64, imageTag string) bool {
	var cb entities.Cluster
	err := r.db.Get(&cb, `SELECT id FROM cluster_builds WHERE cluster_id = ? and image_tag = ?`, clusterId, imageTag)

	if err != nil {
		return false
	}

	return true
}

func (r *ClusterBuildRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM cluster_builds WHERE id = ?`, id)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
