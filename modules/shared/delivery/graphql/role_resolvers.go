package gql

import (
	"errors"
	"strconv"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/validation"
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// RolesListQueryResolver for a list of roles.
func (r *Resolver) RolesListQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := checkPermission(params.Context)
	if !allow {
		return []*domain.Role{}, errors.New("insufficient permission")
	}

	role, err := r.roleUseCase.Fetch(params.Context)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return role, nil
}

// RoleQueryResolver for a single role.
func (r *Resolver) RoleQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := checkPermission(params.Context)
	if !allow {
		return []*domain.Role{}, errors.New("insufficient permission")
	}

	id, err := strconv.ParseInt(params.Args["ID"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	role, err := r.roleUseCase.GetByID(params.Context, id)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return role, nil

}

// RoleCreateResolver creates a new role.
func (r *Resolver) RoleCreateResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := checkPermission(params.Context)
	if !allow {
		return []*domain.Role{}, errors.New("insufficient permission")
	}

	role, err := storeRoleValidation(params)
	if err != nil {
		return nil, err
	}

	if err := r.roleUseCase.Store(params.Context, role); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return nil, nil
}

// RoleUpdateResolver updates the given role.
func (r *Resolver) RoleUpdateResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := checkPermission(params.Context)
	if !allow {
		return []*domain.Role{}, errors.New("insufficient permission")
	}

	role, err := updateRoleValidation(params)
	if err != nil {
		return nil, err
	}

	if err := r.roleUseCase.Update(params.Context, role); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return nil, nil
}

// RoleDeleteResolver deletes the given role.
func (r *Resolver) RoleDeleteResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := checkPermission(params.Context)
	if !allow {
		return []*domain.Role{}, errors.New("insufficient permission")
	}

	id, err := strconv.ParseInt(params.Args["ID"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	err = r.roleUseCase.Delete(params.Context, id)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return nil, nil
}

func storeRoleValidation(params graphql.ResolveParams) (*domain.Role, error) {
	roleParams := params.Args["Role"].(map[string]interface{})

	role := &domain.Role{
		Name:        roleParams["name"].(string),
		Description: roleParams["description"].(string),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := validation.IsAValidSchema(params.Context, role); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return role, nil
}

func updateRoleValidation(params graphql.ResolveParams) (*domain.Role, error) {
	roleParams := params.Args["Role"].(map[string]interface{})

	id, err := strconv.ParseInt(roleParams["id"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	role := &domain.Role{
		ID:          id,
		Name:        roleParams["name"].(string),
		Description: roleParams["description"].(string),
		UpdatedAt:   time.Now(),
	}

	if err := validation.IsAValidSchema(params.Context, role); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return role, nil
}
