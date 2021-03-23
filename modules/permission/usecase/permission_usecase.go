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

func (p *permissionUseCase) Store(ctx context.Context, permission *domain.Permission) error {
	err := p.permissionRepo.Store(ctx, permission)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (p *permissionUseCase) Update(ctx context.Context, permission *domain.Permission) error {
	err := p.permissionRepo.Update(ctx, permission)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (p *permissionUseCase) Delete(ctx context.Context, id int64) error {
	err := p.permissionRepo.Delete(ctx, id)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
