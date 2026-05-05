package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ContextRepository struct {
	db *sqlx.DB
}

type ContextRepositoryInterface interface {
	List(name *string) ([]entities.Context, error)
	GetByID(id int64) (*entities.Context, error)
	Create(ctx *entities.Context) error
	Update(ctx *entities.Context) error
	Delete(id int64) error
}

func NewContextRepository(db *sqlx.DB) *ContextRepository {
	return &ContextRepository{db: db}
}

func (r *ContextRepository) List(name *string) ([]entities.Context, error) {
	query := `SELECT * FROM contexts`
	args := []any{}

	if name != nil {
		query += " WHERE name LIKE ?"
		args = append(args, "%"+*name+"%")
	}

	var ctxs []entities.Context
	if err := r.db.Select(&ctxs, query, args...); err != nil {
		return nil, err
	}

	return ctxs, nil
}

func (r *ContextRepository) GetByID(id int64) (*entities.Context, error) {
	var c entities.Context
	err := r.db.Get(&c, `SELECT * FROM contexts WHERE id = ?`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, err
}

func (r *ContextRepository) ExistsById(id int64) bool {
	var c entities.Context
	err := r.db.Get(&c, `SELECT id FROM contexts WHERE id = ?`, id)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ContextRepository) ExistsByName(name string) bool {
	var c entities.Context
	err := r.db.Get(&c, `SELECT id FROM contexts WHERE name = ?`, name)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ContextRepository) ExistsByNameExceptId(name string, id int64) bool {
	var c entities.Context
	err := r.db.Get(&c, `SELECT id FROM contexts WHERE name = ? and id != ?`, name, id)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}

	return true
}

func (r *ContextRepository) Create(c *entities.Context) error {
	return r.db.QueryRowx(`
		INSERT INTO contexts (name, description)
		VALUES (?, ?)
		RETURNING id, created_at, updated_at`,
		c.Name, c.Description).
		Scan(&c.Id, &c.CreatedAt, &c.UpdatedAt)
}

func (r *ContextRepository) Update(ctx *entities.Context) error {
	res, err := r.db.Exec(`
		UPDATE contexts
		SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, ctx.Name, ctx.Description, ctx.Id)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *ContextRepository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM contexts WHERE id = ?`, id)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
