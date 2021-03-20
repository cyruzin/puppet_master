package gql

import (
	"strconv"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// UsersListQueryResolver for a list of users.
func (r *Resolver) UsersListQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	user, err := r.userUseCase.Fetch(params.Context)
	if err != nil {
		log.Error().Stack().Err(err)
		return nil, err
	}
	return user, nil
}

// UserQueryResolver for a single user.
func (r *Resolver) UserQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	idQuery, isOK := params.Args["id"].(string)

	parsedID, _ := strconv.ParseInt(idQuery, 10, 64)

	if isOK {
		user, err := r.userUseCase.GetByID(params.Context, parsedID)
		if err != nil {
			log.Error().Stack().Err(err)
			return nil, err
		}
		return user, nil
	}
	return &domain.User{}, nil
}

// UserCreateResolver creates a new user.
func (r *Resolver) UserCreateResolver(params graphql.ResolveParams) (interface{}, error) {
	user := &domain.User{
		Name:       params.Args["name"].(string),
		Email:      params.Args["email"].(string),
		Password:   params.Args["password"].(string),
		SuperAdmin: params.Args["superadmin"].(bool),
	}

	err := r.userUseCase.Store(params.Context, user)
	if err != nil {
		log.Error().Stack().Err(err)
		return nil, err
	}
	return nil, nil
}

// UserUpdateResolver updates the given user.
func (r *Resolver) UserUpdateResolver(params graphql.ResolveParams) (interface{}, error) {
	id, err := strconv.ParseInt(params.Args["id"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Err(err)
		return nil, err
	}

	user := &domain.User{
		ID:         id,
		Name:       params.Args["name"].(string),
		Email:      params.Args["email"].(string),
		Password:   params.Args["password"].(string),
		SuperAdmin: params.Args["superadmin"].(bool),
	}

	err = r.userUseCase.Update(params.Context, user)
	if err != nil {
		log.Error().Stack().Err(err)
		return nil, err
	}
	return nil, nil
}
