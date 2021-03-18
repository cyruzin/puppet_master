package postgre

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/jmoiron/sqlx"
)

type postgreRepository struct {
	Conn *sqlx.DB
}

// NewPostgreRoleRepository will create an object that represent the role.Repository interface.
func NewPostgreRoleRepository(Conn *sqlx.DB) domain.RoleRepository {
	return &postgreRepository{Conn}
}

func (p *postgreRepository) Fetch(ctx context.Context) ([]*domain.Role, error) {
	query := `SELECT * FROM roles`

	result := make([]*domain.Role, 0)

	err := p.Conn.SelectContext(ctx, &result, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}

func (p *postgreRepository) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	query := `SELECT * FROM roles WHERE id = ?`

	role := domain.Role{}

	err := p.Conn.GetContext(ctx, &role, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &role, nil
}

func (p *postgreRepository) Store(ctx context.Context, role *domain.Role) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	query := `
	  INSERT INTO roles ( 
		name, 
		description,
		created_at, 
		updated_at
		)
		VALUES (?, ?, ?, ?)
		`

	_, err = p.Conn.ExecContext(
		ctx,
		query,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (p *postgreRepository) Update(ctx context.Context, role *domain.Role) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	query := `
		UPDATE roles
		SET 
		name = ?, 
		description = ?,
		updated_at = ?
		WHERE id = ?
	`

	result, err := p.Conn.ExecContext(
		ctx,
		query,
		role.Name,
		role.Description,
		role.UpdatedAt,
		role.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return errors.New("the resource you requested could not be found")
	}

	tx.Commit()
	return nil
}

func (p *postgreRepository) Delete(ctx context.Context, id int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	query := "DELETE FROM roles WHERE id = ?"

	result, err := p.Conn.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return errors.New("the resource you requested could not be found")
	}

	tx.Commit()
	return nil
}
