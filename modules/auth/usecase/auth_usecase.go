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
	authRepo domain.AuthRepository
}

// NewAuthUsecase will create new an authUsecase object representation
// of domain.AuthUsecase interface.
func NewAuthUsecase(auth domain.AuthRepository) domain.AuthUsecase {
	return &authUseCase{
		authRepo: auth,
	}
}

func (a *authUseCase) Authenticate(ctx context.Context, email, password string) (string, error) {
	hashedPassword, err := a.authRepo.Authenticate(ctx, email)
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	if match := crypto.CheckPasswordHash(password, hashedPassword); !match {
		return "", errors.New("authentication failed")
	}

	token, err := a.GenerateToken()
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	return token, nil
}

func (a *authUseCase) GenerateToken() (string, error) {
	t := jwt.New()
	t.Set(jwt.IssuerKey, viper.GetString(`jwt.issuer`))
	t.Set(jwt.SubjectKey, viper.GetString(`jwt.subject`))
	t.Set(jwt.AudienceKey, viper.GetString(`jwt.audience`))
	t.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*viper.GetDuration(`jwt.expiration`)).Unix())

	payload, err := jwt.Sign(t, jwa.HS256, []byte(viper.GetString(`jwt.secret`)))
	if err != nil {
		log.Error().Stack().Err(err).Msg(err.Error())
		return "", err
	}

	return string(payload), nil
}
