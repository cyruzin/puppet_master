package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/crypto"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type authUseCase struct {
	authRepo       domain.AuthRepository
	permissionRepo domain.PermissionRepository
	roleRepo       domain.RoleRepository
	userRepo       domain.UserRepository
}

// NewAuthUsecase will create new an authUsecase object representation
// of domain.AuthUsecase interface.
func NewAuthUsecase(
	auth domain.AuthRepository,
	permission domain.PermissionRepository,
	role domain.RoleRepository,
	user domain.UserRepository,
) domain.AuthUsecase {
	return &authUseCase{
		authRepo:       auth,
		permissionRepo: permission,
		roleRepo:       role,
		userRepo:       user,
	}
}

func (a *authUseCase) Authenticate(ctx context.Context, email, password string) (string, error) {
	user, err := a.authRepo.Authenticate(ctx, email)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	if match := crypto.CheckPasswordHash(password, user.Password); !match {
		return "", errors.New("authentication failed")
	}

	permissions, err := a.permissionRepo.GetPermissionsByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	roles, err := a.roleRepo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	auth := &domain.Auth{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	if len(permissions) > 1 {
		for _, permission := range permissions {
			auth.Permissions = append(auth.Permissions, permission.Name)
		}
	}

	if len(roles) > 1 {
		for _, role := range roles {
			auth.Roles = append(auth.Roles, role.Name)
		}
	}

	token, err := a.GenerateToken(auth)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	return token, nil
}

func (a *authUseCase) GenerateToken(auth *domain.Auth) (string, error) {
	t := jwt.New()
	t.Set(jwt.IssuerKey, viper.GetString(`jwt.issuer`))
	t.Set(jwt.SubjectKey, viper.GetString(`jwt.subject`))
	t.Set(jwt.AudienceKey, viper.GetString(`jwt.audience`))
	t.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*viper.GetDuration(`jwt.expiration`)).Unix())
	t.Set(`user`, auth)

	payload, err := jwt.Sign(t, jwa.HS256, []byte(viper.GetString(`jwt.secret`)))
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	return string(payload), nil
}
