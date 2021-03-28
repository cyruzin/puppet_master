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
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrFetchError
	}

	return result, nil
}

func (p *postgreRepository) GetByID(ctx context.Context, id int64) (*domain.Permission, error) {
	query := `SELECT * FROM permissions WHERE id = $1`

	permission := domain.Permission{}

	err := p.Conn.GetContext(ctx, &permission, query, id)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrGetByIDError
	}

	return &permission, nil
}

func (p *postgreRepository) Store(ctx context.Context, permission *domain.Permission) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
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
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrStoreError
	}

	return nil
}

func (p *postgreRepository) Update(ctx context.Context, permission *domain.Permission) error {
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

	query := "DELETE FROM permissions WHERE id = $1"

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

	tx.Commit()
	return nil
}

func (p *postgreRepository) GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]*domain.Permission, error) {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrPermissionByID
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := "SELECT * FROM permissions WHERE id = $1"

	permissions := []*domain.Permission{}

	err = p.Conn.GetContext(ctx, &permissions, query, roleID)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrPermissionByID
	}

	return permissions, nil
}

func (p *postgreRepository) GivePermissionToRole(ctx context.Context, permissions []string, roleID int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrAssignPermission
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
	  INSERT INTO permission_role ( 
		 permission_id,
		 role_id
		)
		VALUES ($1, $2)
		`

	for _, permission := range permissions {
		_, err = tx.ExecContext(
			ctx,
			query,
			permission,
			roleID,
		)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrAssignPermission
		}
	}

	return nil
}

func (p *postgreRepository) RemovePermissionToRole(ctx context.Context, permissions []string, roleID int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrRemovePermission
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := "DELETE FROM permission_role WHERE role_id = $1"

	for _, permission := range permissions {
		_, err = tx.ExecContext(
			ctx,
			query,
			permission,
			roleID,
		)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrRemovePermission
		}
	}

	return nil
}

func (p *postgreRepository) SyncPermissionToRole(ctx context.Context, permissions []string, roleID int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrSyncPermission
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = p.RemovePermissionToRole(ctx, permissions, roleID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrSyncPermission
	}

	err = p.GivePermissionToRole(ctx, permissions, roleID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrSyncPermission
	}

	return nil
}

func (p *postgreRepository) GetPermissionsByUserID(ctx context.Context, userID int64) ([]*domain.Permission, error) {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrPermissionByID
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := "SELECT * FROM permissions WHERE id = $1"

	permissions := []*domain.Permission{}

	err = p.Conn.GetContext(ctx, &permissions, query, userID)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrPermissionByID
	}

	return permissions, nil
}

func (p *postgreRepository) GivePermissionToUser(ctx context.Context, permissions []string, userID int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrAssignPermission
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
	  INSERT INTO permission_user ( 
		 permission_id,
		 user_id
		)
		VALUES ($1, $2)
		`

	for _, permission := range permissions {
		_, err = tx.ExecContext(
			ctx,
			query,
			permission,
			userID,
		)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrAssignPermission
		}
	}

	return nil
}

func (p *postgreRepository) RemovePermissionToUser(ctx context.Context, permissions []string, userID int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrRemovePermission
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := "DELETE FROM permission_user WHERE user_id = $1"

	for _, permission := range permissions {
		_, err = tx.ExecContext(
			ctx,
			query,
			permission,
			userID,
		)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrRemovePermission
		}
	}

	return nil
}

func (p *postgreRepository) SyncPermissionToUser(ctx context.Context, permissions []string, userID int64) error {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrSyncPermission
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = p.RemovePermissionToUser(ctx, permissions, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrSyncPermission
	}

	err = p.GivePermissionToUser(ctx, permissions, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrSyncPermission
	}

	return nil
}
