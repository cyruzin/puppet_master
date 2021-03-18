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

// NewPostgrePermissionRepository will create an object that represent
// the permission.Repository interface.
func NewPostgrePermissionRepository(Conn *sqlx.DB) domain.PermissionRepository {
	return &postgreRepository{Conn}
}

func (p *postgreRepository) Fetch(ctx context.Context) ([]*domain.Permission, error) {
	query := `SELECT * FROM permissions`

	result := make([]*domain.Permission, 0)

	err := p.Conn.SelectContext(ctx, &result, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}

func (p *postgreRepository) GetByID(ctx context.Context, id int64) (*domain.Permission, error) {
	query := `SELECT * FROM permissions WHERE id = ?`

	permission := domain.Permission{}

	err := p.Conn.GetContext(ctx, &permission, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &permission, nil
}

func (p *postgreRepository) Store(ctx context.Context, permission *domain.Permission) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	query := `
	  INSERT INTO permissions ( 
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
		permission.Name,
		permission.Description,
		permission.CreatedAt,
		permission.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (p *postgreRepository) Update(ctx context.Context, permission *domain.Permission) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	query := `
		UPDATE permissions
		SET 
		name = ?, 
		description = ?,
		updated_at = ?
		WHERE id = ?
	`

	result, err := p.Conn.ExecContext(
		ctx,
		query,
		permission.Name,
		permission.Description,
		permission.UpdatedAt,
		permission.ID,
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

	query := "DELETE FROM permissions WHERE id = ?"

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
