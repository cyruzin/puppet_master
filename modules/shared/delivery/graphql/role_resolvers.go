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
	if allow := r.authUseCase.Authorize(params.Context, "view role", nil); !allow {
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
	if allow := r.authUseCase.Authorize(params.Context, "view role", nil); !allow {
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
	if allow := r.authUseCase.Authorize(params.Context, "create role", nil); !allow {
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
	if allow := r.authUseCase.Authorize(params.Context, "edit role", nil); !allow {
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
	if allow := r.authUseCase.Authorize(params.Context, "delete role", nil); !allow {
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

func (r *Resolver) RoleGetByUserIDResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.authUseCase.Authorize(params.Context, "get role by user id", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
	}

	userID, ok := params.Args["ID"].(int)
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	role, err := r.roleUseCase.GetRoleByUserID(params.Context, int64(userID))
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *Resolver) RoleAssignResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.authUseCase.Authorize(params.Context, "assign role by user id", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
	}

	roleParams, ok := params.Args["Role"].(map[string]interface{})
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	roleID, ok := roleParams["role_id"].(int)
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	userID, ok := roleParams["user_id"].(int)
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	if err := r.roleUseCase.AssignRoleToUser(params.Context, roleID, int64(userID)); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) RoleSyncResolver(params graphql.ResolveParams) (interface{}, error) {
	if allow := r.authUseCase.Authorize(params.Context, "sync role by user id", nil); !allow {
		log.Error().Err(domain.ErrUnauthorized).Stack().Msg(domain.ErrUnauthorized.Error())
		return nil, domain.ErrUnauthorized
	}

	roleParams, ok := params.Args["Role"].(map[string]interface{})
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	roleID, ok := roleParams["role_id"].(int)
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	userID, ok := roleParams["user_id"].(int)
	if !ok {
		log.Error().Stack().Msg(domain.ErrBadRequest.Error())
		return nil, domain.ErrBadRequest
	}

	if err := r.roleUseCase.SyncRoleToUser(params.Context, roleID, int64(userID)); err != nil {
		return nil, err
	}

	return nil, nil
}
