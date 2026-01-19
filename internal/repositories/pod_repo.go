package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type PodRepository struct {
	db *sqlx.DB
}

type PodRepositoryInterface interface {
	GetByClusterId(clusterId int64) ([]entities.Pod, error)
	Create(pod *entities.Pod) error
	Delete(id int64) error
}

func NewPodRepository(db *sqlx.DB) *PodRepository {
	return &PodRepository{db: db}
}

func (r *PodRepository) GetByClusterId(clusterId int64) ([]entities.Pod, error) {
	query := `SELECT * FROM pods WHERE cluster_id = ? order by id DESC`

	var pods []entities.Pod
	if err := r.db.Select(&pods, query, clusterId); err != nil {
		return nil, err
	}

	return pods, nil
}

func (r *PodRepository) GetById(id int64) (*entities.Pod, error) {
	var pod entities.Pod
	err := r.db.Get(&pod, `SELECT * FROM pods WHERE id = ?`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &pod, err
}

func (r *PodRepository) Create(p *entities.Pod) error {
	return r.db.QueryRow(`
		INSERT INTO pods (cluster_id, container_name, container_id, ip_address)
		VALUES (?, ?, ?, ?)
		RETURNING id, created_at`,
		p.ClusterId,
		p.ContainerName,
		p.ContainerId,
		p.IpAddress,
	).Scan(&p.Id, &p.CreatedAt)
}

func (r *PodRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM pods WHERE id = ?`, id)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
