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

// NewPostgreUserRepository will create an object that represent the user.Repository interface.
func NewPostgreUserRepository(Conn *sqlx.DB) domain.UserRepository {
	return &postgreRepository{Conn}
}

func (p *postgreRepository) Fetch(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT * FROM users`

	result := make([]*domain.User, 0)

	err := p.Conn.SelectContext(ctx, result, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}

func (p *postgreRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT * FROM users WHERE id = ?`

	user := domain.User{}

	err := p.Conn.GetContext(ctx, &user, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &user, nil
}

func (p *postgreRepository) Store(ctx context.Context, user *domain.User) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	query := `
	  INSERT INTO users ( 
		name, 
		email, 
		password,
		superadmin 
		created_at, 
		updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
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
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (p *postgreRepository) Update(ctx context.Context, user *domain.User) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET 
		name = ?, 
		email = ?, 
		password = ?, 
		superadmin = ?, 
		updated_at = ?
		WHERE id = ?
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

	query := "DELETE FROM users WHERE id = ?"

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
