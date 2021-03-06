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

func (a *authUseCase) Authenticate(ctx context.Context, email, password string) (*domain.AuthToken, error) {
	user, err := a.authRepo.Authenticate(ctx, email)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	if match := crypto.CheckPasswordHash(password, user.Password); !match {
		return nil, errors.New("authentication failed")
	}

	role, err := a.roleRepo.GetRoleByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	auth := &domain.Auth{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   role.Name,
	}

	expiration := time.Duration(time.Minute * viper.GetDuration(`jwt.token_expiration`))
	tokenExpiration := time.Now().Add(expiration)

	token, err := a.GenerateToken("user", auth, tokenExpiration)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	auth.Token = token

	refreshToken, err := a.refreshToken(
		"user",
		user.ID,
		time.Now().AddDate(0, 0, viper.GetInt(`jwt.refresh_token_expiration`)),
	)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	payload := &domain.AuthToken{
		Token:        token,
		RefreshToken: refreshToken,
	}

	userCache := &domain.UserCache{}

	if role != nil {
		userCache.ID = user.ID
		userCache.Role = role.Name

		permissions, err := a.permissionRepo.GetPermissionsByRoleName(ctx, role.Name)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return nil, err
		}

		for _, permission := range permissions {
			userCache.Permissions = append(userCache.Permissions, permission.Name)
		}
	}

	if err := a.saveToken(ctx, user.Email, userCache, expiration); err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return payload, nil
}

func (a *authUseCase) Authorize(ctx context.Context, permission string, roles []string) bool {
	user := ctx.Value(domain.ContextKeyID).(map[string]interface{})

	userCache := &domain.UserCache{}

	if err := a.cacheRepo.Get(ctx, user["email"].(string), userCache); err != nil {
		log.Error().Stack().Msg(err.Error())
		return false
	}

	if userCache.Role == "Admin" {
		return true
	}

	for _, currentRole := range roles {
		if currentRole == userCache.Role {
			return true
		}
	}

	for _, currentPermission := range userCache.Permissions {
		if currentPermission == permission {
			return true
		}
	}

	return false
}

func (a *authUseCase) GenerateToken(
	claimKey string,
	claimValue interface{},
	expiration time.Time,
) (string, error) {
	if claimKey == "" || claimValue == nil {
		return "", errors.New("token claim is empty")
	}

	t := jwt.New()
	t.Set(jwt.IssuerKey, viper.GetString(`jwt.issuer`))
	t.Set(jwt.SubjectKey, viper.GetString(`jwt.subject`))
	t.Set(jwt.AudienceKey, viper.GetString(`jwt.audience`))
	t.Set(jwt.ExpirationKey, expiration.Unix())
	t.Set(claimKey, claimValue)

	payload, err := jwt.Sign(t, jwa.HS256, []byte(viper.GetString(`jwt.secret`)))
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	return string(payload), nil
}

func (a *authUseCase) refreshToken(
	claimKey string,
	claimValue interface{},
	expiration time.Time,
) (string, error) {
	if claimKey == "" || claimValue == nil {
		return "", errors.New("refresh token claim is empty")
	}

	t := jwt.New()
	t.Set(jwt.ExpirationKey, expiration.Unix())
	t.Set(claimKey, claimValue)

	payload, err := jwt.Sign(t, jwa.HS256, []byte(viper.GetString(`jwt.secret`)))
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	return string(payload), nil
}

func (a *authUseCase) saveToken(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	if err := a.cacheRepo.Set(ctx, key, value, expiration); err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (a *authUseCase) RefreshToken(ctx context.Context, userID int64) (*domain.AuthToken, error) {
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	role, err := a.roleRepo.GetRoleByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	auth := &domain.Auth{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   role.Name,
	}

	expiration := time.Duration(time.Minute * viper.GetDuration(`jwt.token_expiration`))
	tokenExpiration := time.Now().Add(expiration)

	token, err := a.GenerateToken("user", auth, tokenExpiration)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	auth.Token = token

	refreshToken, err := a.refreshToken(
		"user",
		user.ID,
		time.Now().AddDate(0, 0, viper.GetInt(`jwt.refresh_token_expiration`)),
	)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	payload := &domain.AuthToken{
		Token:        token,
		RefreshToken: refreshToken,
	}

	userCache := &domain.UserCache{}

	if role != nil {
		userCache.ID = user.ID
		userCache.Role = role.Name

		permissions, err := a.permissionRepo.GetPermissionsByRoleName(ctx, role.Name)
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
			return nil, err
		}

		for _, permission := range permissions {
			userCache.Permissions = append(userCache.Permissions, permission.Name)
		}
	}

	if err := a.saveToken(ctx, user.Email, userCache, expiration); err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return nil, err
	}

	return payload, nil
}
