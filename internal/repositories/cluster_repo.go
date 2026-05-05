package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ClusterRepository struct {
	db *sqlx.DB
}

type ClusterRepositoryInterface interface {
	ListByContext(ctxID int64, clusterType *string) ([]entities.Cluster, error)
	Create(cluster *entities.Cluster) error
	Update(cluster *entities.Cluster) error
	Delete(id int64) error
}

func NewClusterRepository(db *sqlx.DB) *ClusterRepository {
	return &ClusterRepository{db: db}
}

func (r *ClusterRepository) ListByContext(ctxId int64, clusterType *string) ([]entities.Cluster, error) {
	query := `SELECT * FROM clusters WHERE context_id = ? order by id DESC`
	args := []any{ctxId}

	if clusterType != nil {
		query += " AND type = ?"
		args = append(args, *clusterType)
	}

	var clusters []entities.Cluster
	if err := r.db.Select(&clusters, query, args...); err != nil {
		return nil, err
	}

	return clusters, nil
}

func (r *ClusterRepository) GetByContextAndId(ctxId int64, id int64) (*entities.Cluster, error) {
	var s entities.Cluster
	err := r.db.Get(&s, `SELECT * FROM clusters WHERE context_id = ? and id = ?`, ctxId, id)

	if err != nil {

		fmt.Println("GetByContextAndId", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &s, err
}

func (r *ClusterRepository) GetById(id int64) (*entities.Cluster, error) {
	var s entities.Cluster
	err := r.db.Get(&s, `SELECT * FROM clusters WHERE id = ?`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &s, err
}

func (r *ClusterRepository) IsExists(ctxId int64, id int64) bool {
	var a entities.Cluster
	err := r.db.Get(&a, `SELECT id FROM clusters WHERE context_id = ? and id = ?`, ctxId, id)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ClusterRepository) ExistsByName(ctxId int64, name string) bool {
	var s entities.Cluster
	err := r.db.Get(&s, `SELECT id FROM clusters WHERE context_id = ? and name = ?`, ctxId, name)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ClusterRepository) ExistsByNameExceptId(ctxId int64, name string, id int64) bool {
	var s entities.Cluster
	err := r.db.Get(&s, `SELECT id FROM clusters WHERE context_id = ? and name = ? and id != ?`, ctxId, name, id)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ClusterRepository) Create(s *entities.Cluster) error {
	return r.db.QueryRow(`
		INSERT INTO clusters (
			context_id, name, ip, port,
			replicas, cpu, memory, path, type, deploy_mode, image, envs
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, status, created_at, updated_at`,
		s.ContextId,
		s.Name,
		s.IP,
		s.Port,
		s.Replicas,
		s.CPU,
		s.Memory,
		s.Path,
		s.Type,
		s.DeployMode,
		s.Image,
		s.Envs,
	).Scan(&s.Id, &s.Status, &s.CreatedAt, &s.UpdatedAt)
}

func (r *ClusterRepository) Update(s *entities.Cluster) error {
	res, err := r.db.NamedExec(`
		UPDATE clusters SET
			name = :name,
			ip = :ip,
			port = :port,
			replicas = :replicas,
			cpu = :cpu,
			memory = :memory,
			path = :path,
			type = :type,
			deploy_mode = :deploy_mode,
			image= :image,
			envs= :envs
		WHERE id = :id
	`, s)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *ClusterRepository) UpdateLatestVersion(s *entities.Cluster) error {
	res, err := r.db.NamedExec(`
		UPDATE clusters SET
			current_image_tag = :current_image_tag,
			current_version = :current_version,
			last_deployed_at = :last_deployed_at
		WHERE id = :id
	`, s)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *ClusterRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM clusters WHERE id = ?`, id)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
