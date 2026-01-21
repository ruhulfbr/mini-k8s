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

func (r *PodRepository) Create(pod *entities.Pod) error {
	return r.db.QueryRow(`
		INSERT INTO pods (cluster_id, container_name, container_id, ip_address, status)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, status, created_at`,
		pod.ClusterId,
		pod.ContainerName,
		pod.ContainerId,
		pod.IpAddress,
		pod.Status,
	).Scan(&pod.Id, &pod.Status, &pod.CreatedAt)
}

func (r *PodRepository) Update(pod *entities.Pod) error {
	res, err := r.db.NamedExec(`
		UPDATE pods SET
			cluster_id     = :cluster_id,
			container_name = :container_name,
			container_id   = :container_id,
			ip_address     = :ip_address,
			status         = :status
		WHERE id = :id
	`, pod)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
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
