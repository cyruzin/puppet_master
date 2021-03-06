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

	permissions := []*domain.Permission{}

	err := p.Conn.SelectContext(ctx, &permissions, query)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrFetchError
	}

	return permissions, nil
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

func (p *postgreRepository) Store(ctx context.Context, permission *domain.Permission) (*domain.Permission, error) {
	tx, err := p.Conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
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
	  INSERT INTO permissions ( 
		name, 
		description,
		created_at, 
		updated_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`

	var lastID int64

	err = p.Conn.GetContext(
		ctx,
		&lastID,
		query,
		permission.Name,
		permission.Description,
		permission.CreatedAt,
		permission.UpdatedAt,
	)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrStoreError
	}

	permission, err = p.GetByID(ctx, lastID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrStoreError
	}

	return permission, nil
}

func (p *postgreRepository) Update(ctx context.Context, permission *domain.Permission) (*domain.Permission, error) {
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
		return nil, domain.ErrUpdateError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrUpdateError
	}

	if rowsAffected == 0 {
		return nil, domain.ErrNotFound
	}

	permission, err = p.GetByID(ctx, permission.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrStoreError
	}

	return permission, nil
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

	query := `SELECT 
							p.id, 
							p.name, 
							p.description, 
							p.created_at, 
							p.updated_at
					 FROM permissions p
					 JOIN permission_role pr ON pr.permission_id = p.id
					 JOIN roles r ON r.id = pr.role_id
					 WHERE r.id = $1
					 GROUP BY p.id`

	permissions := []*domain.Permission{}

	err = p.Conn.SelectContext(ctx, &permissions, query, roleID)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrPermissionByID
	}

	return permissions, nil
}

func (p *postgreRepository) GetPermissionsByRoleName(ctx context.Context, roleName string) ([]*domain.Permission, error) {
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

	query := `SELECT 
							p.id, 
							p.name, 
							p.description, 
							p.created_at, 
							p.updated_at
					 FROM permissions p
					 JOIN permission_role pr ON pr.permission_id = p.id
					 JOIN roles r ON r.id = pr.role_id
					 WHERE r.name = $1
					 GROUP BY p.id`

	permissions := []*domain.Permission{}

	err = p.Conn.SelectContext(ctx, &permissions, query, roleName)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, domain.ErrPermissionByID
	}

	return permissions, nil
}

func (p *postgreRepository) GivePermissionToRole(ctx context.Context, permissions []int, roleID int64) error {
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

	// If the function GivePermissionToRole is called instead of Sync function,
	// previous permissions should be cleaned to avoid duplicates.
	if err := p.RemovePermissionToRole(ctx, permissions, roleID); err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return domain.ErrRemovePermission
	}

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

func (p *postgreRepository) RemovePermissionToRole(ctx context.Context, permissions []int, roleID int64) error {
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

	for i := 0; i <= len(permissions); i++ {
		_, err = tx.ExecContext(
			ctx,
			query,
			roleID,
		)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrRemovePermission
		}
	}

	return nil
}

func (p *postgreRepository) SyncPermissionToRole(ctx context.Context, permissions []int, roleID int64) error {
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

	if len(permissions) > 0 {
		if err := p.RemovePermissionToRole(ctx, permissions, roleID); err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrSyncPermission
		}

		if err := p.GivePermissionToRole(ctx, permissions, roleID); err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrSyncPermission
		}
	} else {
		if err := p.RemovePermissionToRole(ctx, permissions, roleID); err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return domain.ErrSyncPermission
		}
	}

	return nil
}
