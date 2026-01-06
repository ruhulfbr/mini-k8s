package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type PodRepository struct {
	db *sqlx.DB
}

type PodRepositoryInterface interface {
	ListByService(serviceID int64, status *string) ([]entities.Pod, error)
	Create(pod *entities.Pod) error
	Update(pod *entities.Pod) error
	Delete(id int64) error
}

func NewPodRepository(db *sqlx.DB) *PodRepository {
	return &PodRepository{db: db}
}

func (r *PodRepository) ListByService(serviceID int64, status *string) ([]entities.Pod, error) {
	query := `
		SELECT id, application_id, node_id, name, status, created_at
		FROM pods
		WHERE node_id = ?`
	args := []any{serviceID}

	if status != nil {
		query += " AND status = ?"
		args = append(args, *status)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pods []entities.Pod
	for rows.Next() {
		var p entities.Pod
		if err := rows.Scan(
			&p.ID,
			&p.ApplicationID,
			&p.ServiceID,
			&p.Name,
			&p.Status,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		pods = append(pods, p)
	}

	return pods, nil
}

func (r *PodRepository) Create(p *entities.Pod) error {
	return r.db.QueryRow(`
		INSERT INTO pods (application_id, node_id, name, status)
		VALUES (?, ?, ?, ?)
		RETURNING id, created_at`,
		p.ApplicationID,
		p.ServiceID,
		p.Name,
		p.Status,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *PodRepository) Update(p *entities.Pod) error {
	_, err := r.db.Exec(`
		UPDATE pods
		SET name = ?, status = ?
		WHERE id = ?`,
		p.Name,
		p.Status,
		p.ID,
	)
	return err
}

func (r *PodRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM pods WHERE id = ?`, id)
	return err
}
