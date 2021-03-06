package usecase

import (
	"context"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/rs/zerolog/log"
)

type roleUseCase struct {
	permissionRepo domain.PermissionRepository
	roleRepo       domain.RoleRepository
}

// NewRoleUsecase will create new an roleUsecase object representation
// of domain.RoleUsecase interface.
func NewRoleUsecase(
	permission domain.PermissionRepository,
	role domain.RoleRepository,
) domain.RoleUsecase {
	return &roleUseCase{
		permissionRepo: permission,
		roleRepo:       role,
	}
}

func (r *roleUseCase) Fetch(ctx context.Context) ([]*domain.Role, error) {
	roles, err := r.roleRepo.Fetch(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return roles, nil
}

func (r *roleUseCase) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	role, err := r.roleRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return role, nil
}

func (r *roleUseCase) Store(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	newRole, err := r.roleRepo.Store(ctx, role)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return newRole, nil
}

func (r *roleUseCase) Update(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	updatedRole, err := r.roleRepo.Update(ctx, role)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return updatedRole, nil
}

func (r *roleUseCase) Delete(ctx context.Context, id int64) error {
	err := r.roleRepo.Delete(ctx, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (r *roleUseCase) GetRoleByUserID(ctx context.Context, userID int64) (*domain.Role, error) {
	role, err := r.roleRepo.GetRoleByUserID(ctx, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return role, nil
}

func (r *roleUseCase) AssignRoleToUser(ctx context.Context, role int, userID int64) error {
	if err := r.roleRepo.AssignRoleToUser(ctx, role, userID); err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (r *roleUseCase) SyncRoleToUser(ctx context.Context, role int, userID int64) error {
	if err := r.roleRepo.SyncRoleToUser(ctx, role, userID); err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
