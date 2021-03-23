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

	token, err := r.authUseCase.Authenticate(params.Context, user.Email, user.Password)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	auth := &domain.Auth{Token: token}

	return auth, nil
}

func authValidation(params graphql.ResolveParams) (*domain.Auth, error) {
	authParams := params.Args["credentials"].(map[string]interface{})

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
