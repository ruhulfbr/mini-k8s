package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type BuildHistoryRepository struct {
	db *sqlx.DB
}

type BuildHistoryRepositoryInterface interface {
	Create(history *entities.BuildHistory) error
	GetByService(serviceID int64) ([]entities.BuildHistory, error)
	GetRollbackImage(serviceID int64) (*string, error)
	Delete(id int64) error
}

func NewBuildHistoryRepository(db *sqlx.DB) *BuildHistoryRepository {
	return &BuildHistoryRepository{db: db}
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
			&b.Id,
			&b.ApplicationId,
			&b.ServiceId,
			&b.ImageTag,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *BuildHistoryRepository) Create(b *entities.BuildHistory) error {
	return r.db.QueryRow(`
		INSERT INTO build_history (application_id, service_id, version, image_tag)
		VALUES (?, ?, ?, ?)
		RETURNING id, created_at`,
		b.ApplicationId,
		b.ServiceId,
		b.Version,
		b.ImageTag,
	).Scan(&b.Id, &b.CreatedAt)
}

func (r *BuildHistoryRepository) GetById(serviceId int64, id int64) (*entities.BuildHistory, error) {
	var bh entities.BuildHistory
	err := r.db.Get(&bh, `SELECT * FROM build_history WHERE service_id = ? and id = ?`, serviceId, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &bh, err
}

func (r *BuildHistoryRepository) ExistsByVersion(serviceId int64, name string) bool {
	var s entities.Service
	err := r.db.Get(&s, `SELECT id FROM build_history WHERE service_id = ? and name = ?`, serviceId, name)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *BuildHistoryRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM build_history WHERE id = ?`, id)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
