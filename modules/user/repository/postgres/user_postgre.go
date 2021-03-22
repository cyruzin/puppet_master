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

// NewPostgreUserRepository will create an object that represent the user.Repository interface.
func NewPostgreUserRepository(Conn *sqlx.DB) domain.UserRepository {
	return &postgreRepository{Conn}
}

func (p *postgreRepository) Fetch(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT * FROM users`

	result := make([]*domain.User, 0)

	err := p.Conn.SelectContext(ctx, &result, query)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err)
		return nil, domain.ErrFetchError
	}

	return result, nil
}

func (p *postgreRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	user := domain.User{}

	err := p.Conn.GetContext(ctx, &user, query, id)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err)
		return nil, domain.ErrGetByIDError
	}

	return &user, nil
}

func (p *postgreRepository) Store(ctx context.Context, user *domain.User) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err)
		return domain.ErrStoreError
	}

	query := `
	  INSERT INTO users ( 
		name, 
		email, 
		password,
		superadmin,
		created_at, 
		updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		`

	_, err = p.Conn.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.SuperAdmin,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		log.Error().Stack().Err(err)
		tx.Rollback()
		return domain.ErrStoreError
	}

	tx.Commit()
	return nil
}

func (p *postgreRepository) Update(ctx context.Context, user *domain.User) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err)
		return domain.ErrUpdateError
	}

	query := `
		UPDATE users
		SET 
		name = $1, 
		email = $2, 
		password = $3, 
		superadmin = $4, 
		updated_at = $5
		WHERE id = $6
	`

	result, err := p.Conn.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.SuperAdmin,
		user.UpdatedAt,
		user.ID,
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
		return err
	}

	query := "DELETE FROM users WHERE id = $1"

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
