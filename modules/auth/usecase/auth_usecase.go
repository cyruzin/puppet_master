package usecase

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/crypto"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
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
	hashedPassword, err := a.authRepo.Authenticate(ctx, email, password)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		return "", err
	}

	match := crypto.CheckPasswordHash(password, hashedPassword)
	if !match {
		log.Error().Stack().Err(err).Msg("")
		return "", err
	}

	token, err := a.GenerateToken()
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		return "", err
	}

	return token, nil
}

func (a *authUseCase) GenerateToken() (string, error) {
	t := jwt.New()
	t.Set(jwt.IssuerKey, "Puppet Master")
	t.Set(jwt.SubjectKey, "https://github.com/cyruzin/puppet_master")
	t.Set(jwt.AudienceKey, "Auth Services")
	t.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*2).Unix())

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("failed to generate private key: %s\n", err)
		return "", err
	}

	payload, err := jwt.Sign(t, jwa.RS256, privateKey)
	if err != nil {
		log.Error().Stack().Err(err)
		return "", err
	}

	return string(payload), nil
}

func (a *authUseCase) ParseToken(payload string) (interface{}, error) {
	token, err := jwt.ParseString(payload)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("failed to parse payload: %s\n", err)
		return nil, err
	}

	return token, nil
}
