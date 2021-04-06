package gql

import (
	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/validation"
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// AuthQueryResolver autenticathes the given user.
func (r *Resolver) AuthQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	user, err := authValidation(params)
	if err != nil {
		return nil, err
	}

	payload, err := r.authUseCase.Authenticate(params.Context, user.Email, user.Password)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	auth := &domain.AuthToken{Token: payload.Token, RefreshToken: payload.RefreshToken}

	return auth, nil
}

func (r *Resolver) AuthRefreshTokenResolver(params graphql.ResolveParams) (interface{}, error) {
	userID, ok := params.Args["UserID"].(int)
	if !ok {
		log.Error().Stack().Msg(domain.ErrUserID.Error())
		return nil, domain.ErrUserID
	}

	payload, err := r.authUseCase.RefreshToken(params.Context, int64(userID))
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	auth := &domain.AuthToken{Token: payload.Token, RefreshToken: payload.RefreshToken}

	return auth, nil
}

func authValidation(params graphql.ResolveParams) (*domain.Auth, error) {
	authParams, ok := params.Args["Credentials"].(map[string]interface{})
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	auth := &domain.Auth{
		Email:    authParams["email"].(string),
		Password: authParams["password"].(string),
	}

	if err := validation.IsAValidSchema(params.Context, auth); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return auth, nil
}
