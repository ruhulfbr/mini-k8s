package repositories

import (
	"database/sql"
	"errors"

	"github.com/ruhulfbr/mini-k8s/internal/apperrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ApplicationRepository struct {
	db *sql.DB
}

type ApplicationRepositoryInterface interface {
	List(name *string) ([]entities.Application, error)
	GetByID(id int64) (*entities.Application, error)
	Create(app *entities.Application) error
	Update(app *entities.Application) error
	Delete(id int64) error
}

func NewApplicationRepository(db *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) List(name *string) ([]entities.Application, error) {
	query := `SELECT id, name, git_repo, git_branch, created_at, updated_at FROM applications`
	args := []any{}

	if name != nil {
		query += " WHERE name LIKE ?"
		args = append(args, "%"+*name+"%")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []entities.Application
	for rows.Next() {
		var a entities.Application
		if err := rows.Scan(&a.Id, &a.Name, &a.GitRepo, &a.GitBranch, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, nil
}

func (r *ApplicationRepository) GetByID(id int64) (*entities.Application, error) {
	var a entities.Application
	err := r.db.QueryRow(`
		SELECT id, name, git_repo, git_branch, created_at, updated_at
		FROM applications WHERE id = ?`, id).
		Scan(&a.Id, &a.Name, &a.GitRepo, &a.GitBranch, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NotFound
		}
		return nil, err
	}

	return &a, err
}

func (r *ApplicationRepository) ExistsByName(name string) bool {
	var a entities.Application
	err := r.db.QueryRow(`
		SELECT id
		FROM applications WHERE name = ?`, name).
		Scan(&a.Id)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ApplicationRepository) Create(a *entities.Application) error {
	return r.db.QueryRow(`
		INSERT INTO applications (name, git_repo, git_branch)
		VALUES (?, ?)
		RETURNING id, created_at, updated_at`,
		a.Name, a.GitRepo, a.GitBranch).
		Scan(&a.Id, &a.CreatedAt, &a.UpdatedAt)
}

func (r *ApplicationRepository) Update(app *entities.Application) error {
	res, err := r.db.Exec(`
		UPDATE applications
		SET name = ?, git_repo = ?, git_branch = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, app.Name, app.GitRepo, app.GitBranch, app.Id)

	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return apperrors.NotFound
	}

	return nil
}

func (r *ApplicationRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM applications WHERE id = ?`, id)

	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return apperrors.NotFound
	}

	return nil
}
