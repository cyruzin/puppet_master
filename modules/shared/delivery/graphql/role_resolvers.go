package gql

import (
	"strconv"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/validation"
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// RolesListQueryResolver for a list of roles.
func (r *Resolver) RolesListQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.checkPermissions(params.Context, "view role", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
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
	if allow := r.checkPermissions(params.Context, "view role", nil); !allow {
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

	role, err := r.roleUseCase.GetByID(params.Context, parsedID)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return role, nil
}

// RoleCreateResolver creates a new role.
func (r *Resolver) RoleCreateResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.checkPermissions(params.Context, "create role", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
	}

	role, err := storeRoleValidation(params)
	if err != nil {
		return nil, err
	}

	role, err = r.roleUseCase.Store(params.Context, role)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return role, nil
}

// RoleUpdateResolver updates the given role.
func (r *Resolver) RoleUpdateResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.checkPermissions(params.Context, "edit role", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
	}

	role, err := updateRoleValidation(params)
	if err != nil {
		return nil, err
	}

	role, err = r.roleUseCase.Update(params.Context, role)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return role, nil
}

// RoleDeleteResolver deletes the given role.
func (r *Resolver) RoleDeleteResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.checkPermissions(params.Context, "delete role", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
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
	roleParams, ok := params.Args["Role"].(map[string]interface{})
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	parsedPermissions := []int{}

	if roleParams["permissions"] != nil {
		for _, permission := range roleParams["permissions"].([]interface{}) {
			parsedPermissions = append(parsedPermissions, permission.(int))
		}
	}

	role := &domain.Role{
		Name:        roleParams["name"].(string),
		Description: roleParams["description"].(string),
		Permissions: parsedPermissions,
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
	roleParams, ok := params.Args["Role"].(map[string]interface{})
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	id, err := strconv.ParseInt(roleParams["id"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	parsedPermissions := []int{}

	if roleParams["permissions"] != nil {
		for _, permission := range roleParams["permissions"].([]interface{}) {
			parsedPermissions = append(parsedPermissions, permission.(int))
		}
	}

	role := &domain.Role{
		ID:          id,
		Name:        roleParams["name"].(string),
		Description: roleParams["description"].(string),
		Permissions: parsedPermissions,
		UpdatedAt:   time.Now(),
	}

	if err := validation.IsAValidSchema(params.Context, role); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return role, nil
}
