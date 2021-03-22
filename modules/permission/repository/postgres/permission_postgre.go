package postgre

import (
	"context"
	"database/sql"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
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
		log.Error().Stack().Err(err)
		return nil, domain.ErrFetchError
	}

	return result, nil
}

func (p *postgreRepository) GetByID(ctx context.Context, id int64) (*domain.Permission, error) {
	query := `SELECT * FROM permissions WHERE id = $1`

	permission := domain.Permission{}

	err := p.Conn.GetContext(ctx, &permission, query, id)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err)
		return nil, domain.ErrGetByIDError
	}

	return &permission, nil
}

func (p *postgreRepository) Store(ctx context.Context, permission *domain.Permission) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err)
		return domain.ErrStoreError
	}

	query := `
	  INSERT INTO permissions ( 
		name, 
		description,
		created_at, 
		updated_at
		)
		VALUES ($1, $2, $3, $4)
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
		log.Error().Stack().Err(err)
		tx.Rollback()
		return domain.ErrStoreError
	}

	tx.Commit()
	return nil
}

func (p *postgreRepository) Update(ctx context.Context, permission *domain.Permission) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err)
		return domain.ErrUpdateError
	}

	query := `
		UPDATE permissions
		SET 
		name = $1, 
		description = $2,
		updated_at = $3
		WHERE id = $4
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
		log.Error().Stack().Err(err)
		tx.Rollback()
		return domain.ErrUpdateError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Stack().Err(err)
		tx.Rollback()
		return domain.ErrUpdateError
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return domain.ErrNotFound
	}

	tx.Commit()
	return nil
}

func (p *postgreRepository) Delete(ctx context.Context, id int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err)
		return domain.ErrDeleteError
	}

	query := "DELETE FROM permissions WHERE id = $1"

	result, err := p.Conn.ExecContext(ctx, query, id)
	if err != nil {
		log.Error().Stack().Err(err)
		tx.Rollback()
		return domain.ErrDeleteError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Stack().Err(err)
		tx.Rollback()
		return domain.ErrDeleteError
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return domain.ErrNotFound
	}

	tx.Commit()
	return nil
}
