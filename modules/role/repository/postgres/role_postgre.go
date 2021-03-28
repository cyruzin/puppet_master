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
}

// NewPostgreRoleRepository will create an object that represent the role.Repository interface.
func NewPostgreRoleRepository(
	Conn *sqlx.DB,
	permissionRepo domain.PermissionRepository,
) domain.RoleRepository {
	return &postgreRepository{
		Conn,
		permissionRepo,
	}
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

func (p *postgreRepository) Store(ctx context.Context, role *domain.Role) (*domain.Role, error) {
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
	  INSERT INTO roles ( 
		name, 
		description,
		created_at, 
		updated_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`

	var lastID int64

	err = tx.GetContext(
		ctx,
		&lastID,
		query,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrStoreError
	}

	tx.Commit()

	newRole, err := p.GetByID(ctx, lastID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	err = p.permissionRepo.GivePermissionToRole(ctx, role.Permissions, newRole.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return newRole, nil
}

func (p *postgreRepository) Update(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
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
		return nil, domain.ErrUpdateError
	}

	tx.Commit()

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrUpdateError
	}

	if rowsAffected == 0 {
		return nil, domain.ErrNotFound
	}

	err = p.permissionRepo.SyncPermissionToRole(ctx, role.Permissions, role.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	newRole, err := p.GetByID(ctx, role.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return newRole, nil
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

	query := `SELECT 
							r.id, 
							r.name, 
							r.description, 
							r.created_at, 
							r.updated_at
					 FROM roles r
					 JOIN role_user ru ON ru.role_id = r.id
					 JOIN users u ON u.id = ru.role_id
					 WHERE u.id = $1
					 GROUP BY r.id`

	roles := []*domain.Role{}

	err = p.Conn.SelectContext(ctx, &roles, query, userID)
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

	for i := 0; i <= len(roles); i++ {
		_, err = tx.ExecContext(
			ctx,
			query,
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

	if len(roles) > 0 {
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
	} else {
		err = p.RemoveRole(ctx, roles, userID)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrSyncRole
		}
	}

	return nil
}
