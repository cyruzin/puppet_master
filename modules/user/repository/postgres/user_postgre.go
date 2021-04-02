package postgre

import (
	"context"
	"database/sql"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type postgreRepository struct {
	Conn           *sqlx.DB
	permissionRepo domain.PermissionRepository
	roleRepo       domain.RoleRepository
}

// NewPostgreUserRepository will create an object that represent the user.Repository interface.
func NewPostgreUserRepository(
	Conn *sqlx.DB,
	permissionRepo domain.PermissionRepository,
	roleRepo domain.RoleRepository,
) domain.UserRepository {
	return &postgreRepository{
		Conn,
		permissionRepo,
		roleRepo,
	}
}

func (p *postgreRepository) Fetch(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT * FROM users`

	result := make([]*domain.User, 0)

	err := p.Conn.SelectContext(ctx, &result, query)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrFetchError
	}

	return result, nil
}

func (p *postgreRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	user := domain.User{}

	err := p.Conn.GetContext(ctx, &user, query, id)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrGetByIDError
	}

	return &user, nil
}

func (p *postgreRepository) Store(ctx context.Context, user *domain.User) (*domain.User, error) {
	tx, err := p.Conn.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrStoreError
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
	  INSERT INTO users ( 
		name, 
		email, 
		password,
		created_at, 
		updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`

	var lastID int64

	err = tx.GetContext(
		ctx,
		&lastID,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrStoreError
	}

	tx.Commit()

	newUser, err := p.GetByID(ctx, lastID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrStoreError
	}

	if user.Role > 0 {
		err = p.roleRepo.AssignRole(ctx, user.Role, lastID)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return nil, err
		}
	}

	return newUser, nil
}

func (p *postgreRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	tx, err := p.Conn.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrUpdateError
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
		UPDATE users
		SET 
		name = $1, 
		email = $2, 
		updated_at = $3
		WHERE id = $4
	`
	result, err := tx.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrUpdateError
	}

	tx.Commit()

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrDeleteError
	}

	if rowsAffected == 0 {
		return nil, domain.ErrNotFound
	}

	updatedUser, err := p.GetByID(ctx, user.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrStoreError
	}

	err = p.roleRepo.SyncRole(ctx, user.Role, updatedUser.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return updatedUser, nil
}

func (p *postgreRepository) Delete(ctx context.Context, id int64) error {
	tx, err := p.Conn.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := "DELETE FROM users WHERE id = $1"

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrDeleteError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrDeleteError
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
