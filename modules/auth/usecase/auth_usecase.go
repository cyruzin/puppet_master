package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/crypto"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type authUseCase struct {
	authRepo       domain.AuthRepository
	cacheRepo      domain.CacheRepository
	permissionRepo domain.PermissionRepository
	roleRepo       domain.RoleRepository
	userRepo       domain.UserRepository
}

// NewAuthUsecase will create new an authUsecase object representation
// of domain.AuthUsecase interface.
func NewAuthUsecase(
	auth domain.AuthRepository,
	cache domain.CacheRepository,
	permission domain.PermissionRepository,
	role domain.RoleRepository,
	user domain.UserRepository,
) domain.AuthUsecase {
	return &authUseCase{
		authRepo:       auth,
		cacheRepo:      cache,
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

	roles, err := a.roleRepo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	auth := &domain.Auth{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}

	if len(roles) >= 1 {
		for _, role := range roles {
			auth.Roles = append(auth.Roles, role.Name)
		}
	}

	tokenUUID, err := uuid.NewUUID()
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	auth.TokenUUID = tokenUUID
	auth.Revoked = false

	expiration := time.Now().Add(time.Hour * viper.GetDuration(`jwt.expiration`))

	token, err := a.GenerateToken(auth, expiration)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	auth.Token = token

	a.saveToken(ctx, auth, expiration)

	return token, nil
}

func (a *authUseCase) GenerateToken(auth *domain.Auth, expiration time.Time) (string, error) {
	t := jwt.New()
	t.Set(jwt.IssuerKey, viper.GetString(`jwt.issuer`))
	t.Set(jwt.SubjectKey, viper.GetString(`jwt.subject`))
	t.Set(jwt.AudienceKey, viper.GetString(`jwt.audience`))
	t.Set(jwt.ExpirationKey, expiration.Unix())
	t.Set(`user`, auth)

	payload, err := jwt.Sign(t, jwa.HS256, []byte(viper.GetString(`jwt.secret`)))
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	return string(payload), nil
}

func (a *authUseCase) saveToken(ctx context.Context, auth *domain.Auth, expiration time.Time) error {
	if err := a.cacheRepo.Set(ctx, auth.Email, auth, time.Duration(expiration.Unix())); err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

// func (a *authUseCase) getToken(ctx context.Context, key string, auth *domain.Auth) error {
// 	err := a.cacheRepo.Get(ctx, key, auth)
// 	if err != nil {
// 		log.Error().Stack().Err(err).Msg(err.Error())
// 		return err
// 	}

// 	return nil
// }
