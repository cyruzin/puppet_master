package gql

import (
	"strconv"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/crypto"
	"github.com/cyruzin/puppet_master/pkg/validation"
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// UsersListQueryResolver for a list of users.
func (r *Resolver) UsersListQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	users, err := r.userUseCase.Fetch(params.Context)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return users, nil
}

// UserQueryResolver for a single user.
func (r *Resolver) UserQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	idQuery, isOK := params.Args["ID"].(string)

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
	user, err := storeUserValidation(params)
	if err != nil {
		return nil, err
	}

	if err := r.userUseCase.Store(params.Context, user); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return nil, nil
}

// UserUpdateResolver updates the given user.
func (r *Resolver) UserUpdateResolver(params graphql.ResolveParams) (interface{}, error) {
	user, err := updateUserValidation(params)
	if err != nil {
		return nil, err
	}

	if err := r.userUseCase.Update(params.Context, user); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return nil, nil
}

// UserDeleteResolver deletes the given user.
func (r *Resolver) UserDeleteResolver(params graphql.ResolveParams) (interface{}, error) {
	id, err := strconv.ParseInt(params.Args["ID"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	err = r.userUseCase.Delete(params.Context, id)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return nil, nil
}

func storeUserValidation(params graphql.ResolveParams) (*domain.User, error) {
	userParams := params.Args["User"].(map[string]interface{})

	password := userParams["password"].(string)

	user := &domain.User{
		Name:       userParams["name"].(string),
		Email:      userParams["email"].(string),
		Password:   password,
		SuperAdmin: userParams["superadmin"].(bool),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := validation.IsAValidSchema(params.Context, user); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	hashedPassword, err := crypto.HashPassword(password, 6)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	user.Password = hashedPassword

	return user, nil
}

func updateUserValidation(params graphql.ResolveParams) (*domain.User, error) {
	userParams := params.Args["User"].(map[string]interface{})

	id, err := strconv.ParseInt(userParams["id"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	user := &domain.User{
		ID:         id,
		Name:       userParams["name"].(string),
		Email:      userParams["email"].(string),
		SuperAdmin: userParams["superadmin"].(bool),
		UpdatedAt:  time.Now(),
	}

	if err := validation.IsAValidSchema(params.Context, user); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return user, nil
}
