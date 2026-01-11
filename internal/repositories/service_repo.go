package repositories

import (
	"database/sql"
	"errors"
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
	query := `SELECT * FROM services WHERE application_id = ? order by id DESC`
	args := []any{appID}

	if serviceType != nil {
		query += " AND type = ?"
		args = append(args, *serviceType)
	}

	var services []entities.Service
	if err := r.db.Select(&services, query, args...); err != nil {
		return nil, err
	}

	return services, nil
}

func (r *ServiceRepository) GetById(appId int64, id int64) (*entities.Service, error) {
	var s entities.Service
	err := r.db.Get(&s, `SELECT * FROM services WHERE application_id = ? and id = ?`, appId, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &s, err
}

func (r *ServiceRepository) IsExists(appId int64, id int64) bool {
	var a entities.Service
	err := r.db.Get(&a, `SELECT id FROM services WHERE application_id = ? and id = ?`, appId, id)

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
	err := r.db.Get(&s, `SELECT id FROM services WHERE application_id = ? and name = ? and id != ?`, appId, name, id)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ServiceRepository) Create(s *entities.Service) error {
	return r.db.QueryRow(`
		INSERT INTO services (
			application_id, name, ip, port,
			replicas, cpu, memory, path, type, current_image_tag
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, status, created_at, updated_at`,
		s.ApplicationId,
		s.Name,
		s.IP,
		s.Port,
		s.Replicas,
		s.CPU,
		s.Memory,
		s.Path,
		s.Type,
		s.CurrentImageTag,
	).Scan(&s.Id, &s.Status, &s.CreatedAt, &s.UpdatedAt)
}

func (r *ServiceRepository) Update(s *entities.Service) error {
	res, err := r.db.NamedExec(`
		UPDATE services SET
			name = :name,
			ip = :ip,
			port = :port,
			image = :image,
			replicas = :replicas,
			cpu = :cpu,
			memory = :memory,
			path = :path,
			type = :type,
			current_image_tag= :current_image_tag
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

func (r *ServiceRepository) UpdateLastDeployed(id int64, imageTag string, version string) error {
	res, err := r.db.Exec(`
		UPDATE services
		SET last_deployed_at = ?, current_image_tag = ?, version = ? 
		WHERE id = ?`,
		time.Now(), id, imageTag, version,
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
