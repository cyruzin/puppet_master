package gql

import (
	"strconv"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/crypto"
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// UsersListQueryResolver for a list of users.
func (r *Resolver) UsersListQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	user, err := r.userUseCase.Fetch(params.Context)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
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
			log.Error().Stack().Msg(err.Error())
			return nil, err
		}
		return user, nil
	}
	return &domain.User{}, nil
}

// UserCreateResolver creates a new user.
func (r *Resolver) UserCreateResolver(params graphql.ResolveParams) (interface{}, error) {
	userParams := params.Args["user"].(map[string]interface{})

	superadmin := false

	if userParams["superadmin"].(bool) {
		superadmin = true
	}

	password, err := crypto.HashPassword(userParams["password"].(string), 6)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	user := &domain.User{
		Name:       userParams["name"].(string),
		Email:      userParams["email"].(string),
		Password:   password,
		SuperAdmin: superadmin,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = r.userUseCase.Store(params.Context, user)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return nil, nil
}

// UserUpdateResolver updates the given user.
func (r *Resolver) UserUpdateResolver(params graphql.ResolveParams) (interface{}, error) {
	userParams := params.Args["user"].(map[string]interface{})

	id, err := strconv.ParseInt(userParams["id"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	superadmin := false

	if userParams["superadmin"].(bool) {
		superadmin = true
	}

	password, err := crypto.HashPassword(userParams["password"].(string), 6)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	user := &domain.User{
		ID:         id,
		Name:       userParams["name"].(string),
		Email:      userParams["email"].(string),
		Password:   password,
		SuperAdmin: superadmin,
		UpdatedAt:  time.Now(),
	}

	err = r.userUseCase.Update(params.Context, user)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return nil, nil
}
