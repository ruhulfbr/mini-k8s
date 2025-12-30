package repositories

import (
	"database/sql"

	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type BuildHistoryRepository struct {
	db *sql.DB
}

type BuildHistoryRepositoryInterface interface {
	Create(history *entities.BuildHistory) error
	GetByApplication(appID int64) ([]entities.BuildHistory, error)
	GetByService(serviceID int64) ([]entities.BuildHistory, error)
	GetRollbackImage(serviceID int64) (*string, error)
	Delete(id int64) error
}

func NewBuildHistoryRepository(db *sql.DB) *BuildHistoryRepository {
	return &BuildHistoryRepository{db: db}
}

func (r *BuildHistoryRepository) Create(b *entities.BuildHistory) error {
	return r.db.QueryRow(`
		INSERT INTO build_history (application_id, node_id, tag)
		VALUES (?, ?, ?)
		RETURNING id, created_at`,
		b.ApplicationID,
		b.ServiceID,
		b.Tag,
	).Scan(&b.ID, &b.CreatedAt)
}

func (r *BuildHistoryRepository) GetByApplication(appID int64) ([]entities.BuildHistory, error) {
	rows, err := r.db.Query(`
		SELECT id, application_id, node_id, tag, created_at
		FROM build_history
		WHERE application_id = ?
		ORDER BY created_at DESC`,
		appID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.BuildHistory
	for rows.Next() {
		var b entities.BuildHistory
		if err := rows.Scan(
			&b.ID,
			&b.ApplicationID,
			&b.ServiceID,
			&b.Tag,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *BuildHistoryRepository) GetByService(serviceID int64) ([]entities.BuildHistory, error) {
	rows, err := r.db.Query(`
		SELECT id, application_id, node_id, tag, created_at
		FROM build_history
		WHERE node_id = ?
		ORDER BY created_at DESC`,
		serviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.BuildHistory
	for rows.Next() {
		var b entities.BuildHistory
		if err := rows.Scan(
			&b.ID,
			&b.ApplicationID,
			&b.ServiceID,
			&b.Tag,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *BuildHistoryRepository) GetRollbackImage(serviceID int64) (*string, error) {
	var tag string
	err := r.db.QueryRow(`
		SELECT tag FROM build_history
		WHERE node_id = ?
		ORDER BY created_at DESC
		LIMIT 1 OFFSET 1`,
		serviceID,
	).Scan(&tag)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &tag, err
}

func (r *BuildHistoryRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM build_history WHERE id = ?`, id)
	return err
}
