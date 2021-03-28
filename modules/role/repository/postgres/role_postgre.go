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

// NewPostgreRoleRepository will create an object that represent the role.Repository interface.
func NewPostgreRoleRepository(Conn *sqlx.DB) domain.RoleRepository {
	return &postgreRepository{Conn}
}

func (p *postgreRepository) Fetch(ctx context.Context) ([]*domain.Role, error) {
	query := `SELECT * FROM roles`

	result := make([]*domain.Role, 0)

	err := p.Conn.SelectContext(ctx, &result, query)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrFetchError
	}

	return result, nil
}

func (p *postgreRepository) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	query := `SELECT * FROM roles WHERE id = $1`

	role := domain.Role{}

	err := p.Conn.GetContext(ctx, &role, query, id)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrGetByIDError
	}

	return &role, nil
}

func (p *postgreRepository) Store(ctx context.Context, role *domain.Role) error {
	tx, err := p.Conn.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrStoreError
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
	  INSERT INTO roles ( 
		name, 
		description,
		created_at, 
		updated_at
		)
		VALUES ($1, $2, $3, $4)
		`

	_, err = tx.ExecContext(
		ctx,
		query,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrStoreError
	}

	return nil
}

func (p *postgreRepository) Update(ctx context.Context, role *domain.Role) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrUpdateError
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
		UPDATE roles
		SET 
		name = $1, 
		description = $2,
		updated_at = $3
		WHERE id = $4
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
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrUpdateError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrUpdateError
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (p *postgreRepository) Delete(ctx context.Context, id int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrDeleteError
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := "DELETE FROM roles WHERE id = $1"

	result, err := p.Conn.ExecContext(ctx, query, id)
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

func (p *postgreRepository) GetRolesByUserID(ctx context.Context, userID int64) ([]*domain.Role, error) {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrRoleByID
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `SELECT * FROM roles WHERE id = $1`

	roles := []*domain.Role{}

	err = p.Conn.GetContext(ctx, &roles, query, userID)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrRoleByID
	}

	return roles, nil
}

func (p *postgreRepository) AssignRole(ctx context.Context, roles []int, userID int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrAssignRole
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
	  INSERT INTO role_user ( 
		 role_id,
		 user_id
		)
		VALUES ($1, $2)
		`

	for _, role := range roles {
		_, err = tx.ExecContext(
			ctx,
			query,
			role,
			userID,
		)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrAssignRole
		}
	}

	return nil
}

func (p *postgreRepository) RemoveRole(ctx context.Context, roles []int, userID int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrRemoveRole
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := "DELETE FROM role_user WHERE user_id = $1"

	for _, role := range roles {
		_, err = tx.ExecContext(
			ctx,
			query,
			role,
			userID,
		)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrRemoveRole
		}
	}

	return nil
}

func (p *postgreRepository) SyncRole(ctx context.Context, roles []int, userID int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrSyncRole
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = p.RemoveRole(ctx, roles, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrSyncRole
	}

	err = p.AssignRole(ctx, roles, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrSyncRole
	}

	return nil
}
