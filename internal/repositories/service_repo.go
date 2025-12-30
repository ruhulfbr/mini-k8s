package repositories

import (
	"database/sql"
	"time"

	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ServiceRepository struct {
	db *sql.DB
}

type ServiceRepositoryInterface interface {
	ListByApplication(appID int64, status *string) ([]entities.Service, error)
	Create(service *entities.Service) error
	Update(service *entities.Service) error
	Delete(id int64) error
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) ListByApplication(appID int64, serviceType *string) ([]entities.Service, error) {
	query := `
		SELECT id, application_id, name, ip, port, image_tag,
		       context_path, replicas, resources, path, type, last_build_at
		FROM nodes
		WHERE application_id = ?`
	args := []any{appID}

	if serviceType != nil {
		query += " AND type = ?"
		args = append(args, *serviceType)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []entities.Service
	for rows.Next() {
		var s entities.Service
		var port sql.NullInt64
		var lastBuild sql.NullTime

		if err := rows.Scan(
			&s.ID,
			&s.ApplicationID,
			&s.Name,
			&s.IP,
			&port,
			&s.ImageTag,
			&s.ContextPath,
			&s.Replicas,
			&s.Resources,
			&s.Path,
			&s.Type,
			&lastBuild,
		); err != nil {
			return nil, err
		}

		if port.Valid {
			p := int(port.Int64)
			s.Port = &p
		}

		if lastBuild.Valid {
			t := lastBuild.Time
			s.LastBuildAt = &t
		}

		services = append(services, s)
	}

	return services, nil
}

func (r *ServiceRepository) Create(s *entities.Service) error {
	var lastBuild sql.NullTime
	if s.LastBuildAt != nil {
		lastBuild = sql.NullTime{Time: *s.LastBuildAt, Valid: true}
	}

	return r.db.QueryRow(`
		INSERT INTO nodes (
			application_id, name, ip, port, image_tag,
			context_path, replicas, resources, path, type, last_build_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`,
		s.ApplicationID,
		s.Name,
		s.IP,
		s.Port,
		s.ImageTag,
		s.ContextPath,
		s.Replicas,
		s.Resources,
		s.Path,
		s.Type,
		lastBuild,
	).Scan(&s.ID)
}

func (r *ServiceRepository) Update(s *entities.Service) error {
	_, err := r.db.Exec(`
		UPDATE nodes SET
			name = ?,
			ip = ?,
			port = ?,
			image_tag = ?,
			context_path = ?,
			replicas = ?,
			resources = ?,
			path = ?,
			type = ?,
			last_build_at = ?
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
		s.LastBuildAt,
		s.ID,
	)
	return err
}

func (r *ServiceRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM nodes WHERE id = ?`, id)
	return err
}

func (r *ServiceRepository) TouchLastBuild(serviceID int64) error {
	_, err := r.db.Exec(`
		UPDATE nodes
		SET last_build_at = ?
		WHERE id = ?`,
		time.Now(), serviceID,
	)
	return err
}
