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
	allow := r.checkPermission(params.Context)
	if !allow {
		return []*domain.User{}, domain.ErrUnauthorized
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
	allow := r.checkPermission(params.Context)
	if !allow {
		return []*domain.User{}, domain.ErrUnauthorized
	}

	id, err := strconv.ParseInt(params.Args["ID"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	user, err := r.userUseCase.GetByID(params.Context, id)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return user, nil

}

// UserCreateResolver creates a new user.
func (r *Resolver) UserCreateResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := r.checkPermission(params.Context)
	if !allow {
		return []*domain.User{}, domain.ErrUnauthorized
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
	allow := r.checkPermission(params.Context)
	if !allow {
		return []*domain.User{}, domain.ErrUnauthorized
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
	allow := r.checkPermission(params.Context)
	if !allow {
		return []*domain.User{}, domain.ErrUnauthorized
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
	userParams := params.Args["User"].(map[string]interface{})

	parsedRoles := []int{}

	if userParams["roles"] != nil {
		for _, role := range userParams["roles"].([]interface{}) {
			parsedRoles = append(parsedRoles, role.(int))
		}
	}

	parsedPermissions := []int{}

	if userParams["permissions"] != nil {
		for _, role := range userParams["permissions"].([]interface{}) {
			parsedPermissions = append(parsedPermissions, role.(int))
		}
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
		Name:        userParams["name"].(string),
		Email:       userParams["email"].(string),
		Roles:       parsedRoles,
		Permissions: parsedPermissions,
		Password:    password,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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

	parsedRoles := []int{}

	if userParams["roles"] != nil {
		for _, role := range userParams["roles"].([]interface{}) {
			parsedRoles = append(parsedRoles, role.(int))
		}
	}

	parsedPermissions := []int{}

	if userParams["permissions"] != nil {
		for _, role := range userParams["permissions"].([]interface{}) {
			parsedPermissions = append(parsedPermissions, role.(int))
		}
	}

	user := &domain.User{
		ID:          id,
		Name:        userParams["name"].(string),
		Email:       userParams["email"].(string),
		Roles:       parsedRoles,
		Permissions: parsedPermissions,
		UpdatedAt:   time.Now(),
	}

	if err := validation.IsAValidSchema(params.Context, user); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return user, nil
}
