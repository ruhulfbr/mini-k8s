package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ServiceRepository struct {
	db *sqlx.DB
}

type ServiceRepositoryInterface interface {
	ListByApplication(appID int64, status *string) ([]entities.Service, error)
	Create(service *entities.Service) error
	Update(service *entities.Service) error
	Delete(id int64) error
}

func NewServiceRepository(db *sqlx.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) ListByApplication(appID int64, serviceType *string) ([]entities.Service, error) {
	query := `SELECT * FROM services WHERE application_id = ?`
	args := []any{appID}

	if serviceType != nil {
		query += " AND type = ?"
		args = append(args, *serviceType)
	}

	var services []entities.Service
	if err := r.db.Select(&services, query, args...); err != nil {

		fmt.Println(err)

		return nil, err
	}

	return services, nil
}

func (r *ServiceRepository) ExistsById(id int64) bool {
	var a entities.Service
	err := r.db.Get(&a, `SELECT id FROM services WHERE id = ?`, id)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ServiceRepository) ExistsByName(appId int64, name string) bool {
	var s entities.Service
	err := r.db.Get(&s, `SELECT id FROM services WHERE application_id = ? and name = ?`, appId, name)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ServiceRepository) ExistsByNameExceptId(appId int64, name string, id int64) bool {
	var s entities.Service
	err := r.db.Get(&s, `SELECT id FROM services WHERE application_i = ? and name = ? and id != ?`, appId, name, id)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ServiceRepository) Create(s *entities.Service) error {
	return r.db.QueryRow(`
		INSERT INTO services (
			application_id, name, ip, port, image_tag,
			context_path, replicas, resources, path, type
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at`,
		s.ApplicationId,
		s.Name,
		s.IP,
		s.Port,
		s.ImageTag,
		s.ContextPath,
		s.Replicas,
		s.Resources,
		s.Path,
		s.Type,
	).Scan(&s.Id, &s.CreatedAt, &s.UpdatedAt)
}

func (r *ServiceRepository) Update(s *entities.Service) error {
	res, err := r.db.Exec(`
		UPDATE services SET
			name = ?,
			ip = ?,
			port = ?,
			image_tag = ?,
			context_path = ?,
			replicas = ?,
			resources = ?,
			path = ?,
			type = ?,
		WHERE id = ?`,
		s.Name,
		s.IP,
		s.Port,
		s.ImageTag,
		s.ContextPath,
		s.Replicas,
		s.Resources,
		s.Path,
		s.Type,
		s.Id,
	)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *ServiceRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM services WHERE id = ?`, id)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *ServiceRepository) TouchLastBuild(id int64) error {
	res, err := r.db.Exec(`
		UPDATE services
		SET last_build_at = ?
		WHERE id = ?`,
		time.Now(), id,
	)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
