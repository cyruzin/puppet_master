package usecase

import (
	"context"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/rs/zerolog/log"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase will create new an articleUsecase object representation
// of domain.UserUsecase interface.
func NewUserUsecase(user domain.UserRepository) domain.UserUsecase {
	return &userUseCase{
		userRepo: user,
	}
}

func (u *userUseCase) Fetch(ctx context.Context) ([]*domain.User, error) {
	users, err := u.userRepo.Fetch(ctx)
	if err != nil {
		log.Error().Stack().Err(err)
		return nil, err
	}

	return users, nil
}

func (u *userUseCase) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Stack().Err(err)
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Store(ctx context.Context, user *domain.User) error {
	err := u.userRepo.Store(ctx, user)
	if err != nil {
		log.Error().Stack().Err(err)
		return err
	}

	return nil
}

func (u *userUseCase) Update(ctx context.Context, user *domain.User) error {
	err := u.userRepo.Update(ctx, user)
	if err != nil {
		log.Error().Stack().Err(err)
		return err
	}

	return nil
}

func (u *userUseCase) Delete(ctx context.Context, id int64) error {
	err := u.userRepo.Delete(ctx, id)
	if err != nil {
		log.Error().Stack().Err(err)
		return err
	}

	return nil
}
