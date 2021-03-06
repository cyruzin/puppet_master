package usecase

import (
	"context"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/rs/zerolog/log"
)

type permissionUseCase struct {
	permissionRepo domain.PermissionRepository
}

// NewPermissionUsecase will create new an permissionUsecase object
// representation of domain.PermissionUsecase interface.
func NewPermissionUsecase(permission domain.PermissionRepository) domain.PermissionUsecase {
	return &permissionUseCase{
		permissionRepo: permission,
	}
}

func (p *permissionUseCase) Fetch(ctx context.Context) ([]*domain.Permission, error) {
	permissions, err := p.permissionRepo.Fetch(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return permissions, nil
}

func (p *permissionUseCase) GetByID(ctx context.Context, id int64) (*domain.Permission, error) {
	permission, err := p.permissionRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return permission, nil
}

func (p *permissionUseCase) Store(ctx context.Context, permission *domain.Permission) (*domain.Permission, error) {
	permission, err := p.permissionRepo.Store(ctx, permission)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return permission, nil
}

func (p *permissionUseCase) Update(ctx context.Context, permission *domain.Permission) (*domain.Permission, error) {
	permission, err := p.permissionRepo.Update(ctx, permission)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return permission, nil
}

func (p *permissionUseCase) Delete(ctx context.Context, id int64) error {
	err := p.permissionRepo.Delete(ctx, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (p *permissionUseCase) GetPermissionsByRoleID(ctx context.Context, roleID int64) ([]*domain.Permission, error) {
	permissions, err := p.permissionRepo.GetPermissionsByRoleID(ctx, roleID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return permissions, nil
}

func (p *permissionUseCase) GetPermissionsByRoleName(ctx context.Context, roleName string) ([]*domain.Permission, error) {
	permissions, err := p.permissionRepo.GetPermissionsByRoleName(ctx, roleName)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return permissions, nil
}

func (p *permissionUseCase) GivePermissionToRole(ctx context.Context, permissions []int, roleID int64) error {
	if err := p.permissionRepo.GivePermissionToRole(ctx, permissions, roleID); err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (p *permissionUseCase) SyncPermissionToRole(ctx context.Context, permissions []int, roleID int64) error {
	if err := p.permissionRepo.SyncPermissionToRole(ctx, permissions, roleID); err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
