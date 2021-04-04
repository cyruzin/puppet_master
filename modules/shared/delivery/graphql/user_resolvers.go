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
	if allow := r.authUseCase.Authorize(params.Context, "view user", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
	}

	users, err := r.userUseCase.Fetch(params.Context)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return users, nil
}

// UserQueryResolver for a single user.
func (r *Resolver) UserQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.authUseCase.Authorize(params.Context, "view user", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
	}

	id, ok := params.Args["ID"].(string)
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	user, err := r.userUseCase.GetByID(params.Context, parsedID)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return user, nil
}

// UserCreateResolver creates a new user.
func (r *Resolver) UserCreateResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.authUseCase.Authorize(params.Context, "create user", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
	}

	user, err := storeUserValidation(params)
	if err != nil {
		return nil, err
	}

	user, err = r.userUseCase.Store(params.Context, user)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return user, nil
}

// UserUpdateResolver updates the given user.
func (r *Resolver) UserUpdateResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.authUseCase.Authorize(params.Context, "delete user", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
	}

	user, err := updateUserValidation(params)
	if err != nil {
		return nil, err
	}

	user, err = r.userUseCase.Update(params.Context, user)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return user, nil
}

// UserDeleteResolver deletes the given user.
func (r *Resolver) UserDeleteResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.authUseCase.Authorize(params.Context, "edit user", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
	}

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
	userParams, ok := params.Args["User"].(map[string]interface{})
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	password := ""

	if userParams["password"] != nil {
		password = userParams["password"].(string)
	}

	err := validation.IsAValidField(params.Context, password, "password", "required,gte=8")
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	user := &domain.User{
		Name:      userParams["name"].(string),
		Email:     userParams["email"].(string),
		Role:      userParams["role"].(int),
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
	userParams, ok := params.Args["User"].(map[string]interface{})
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	id, err := strconv.ParseInt(userParams["id"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	user := &domain.User{
		ID:        id,
		Name:      userParams["name"].(string),
		Email:     userParams["email"].(string),
		Role:      userParams["role"].(int),
		UpdatedAt: time.Now(),
	}

	if err := validation.IsAValidSchema(params.Context, user); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return user, nil
}
