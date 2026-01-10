package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ApplicationRepository struct {
	db *sqlx.DB
}

type ApplicationRepositoryInterface interface {
	List(name *string) ([]entities.Application, error)
	GetByID(id int64) (*entities.Application, error)
	Create(app *entities.Application) error
	Update(app *entities.Application) error
	Delete(id int64) error
}

func NewApplicationRepository(db *sqlx.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) List(name *string) ([]entities.Application, error) {
	query := `SELECT * FROM applications`
	args := []any{}

	if name != nil {
		query += " WHERE name LIKE ?"
		args = append(args, "%"+*name+"%")
	}

	var apps []entities.Application
	if err := r.db.Select(&apps, query, args...); err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *ApplicationRepository) GetByID(id int64) (*entities.Application, error) {
	var a entities.Application
	err := r.db.Get(&a, `SELECT * FROM applications WHERE id = ?`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &a, err
}

func (r *ApplicationRepository) ExistsById(id int64) bool {
	var a entities.Application
	err := r.db.Get(&a, `SELECT id FROM applications WHERE id = ?`, id)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ApplicationRepository) ExistsByName(name string) bool {
	var a entities.Application
	err := r.db.Get(&a, `SELECT id FROM applications WHERE name = ?`, name)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ApplicationRepository) ExistsByNameExceptId(name string, id int64) bool {
	var a entities.Application
	err := r.db.Get(&a, `SELECT id FROM applications WHERE name = ? and id != ?`, name, id)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ApplicationRepository) Create(a *entities.Application) error {
	return r.db.QueryRowx(`
		INSERT INTO applications (name, description)
		VALUES (?, ?)
		RETURNING id, created_at, updated_at`,
		a.Name, a.Description).
		Scan(&a.Id, &a.CreatedAt, &a.UpdatedAt)
}

func (r *ApplicationRepository) Update(app *entities.Application) error {
	res, err := r.db.Exec(`
		UPDATE applications
		SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, app.Name, app.Description, app.Id)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *ApplicationRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM applications WHERE id = ?`, id)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
