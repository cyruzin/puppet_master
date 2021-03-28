package gql

import (
	"strconv"
	"time"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/cyruzin/puppet_master/pkg/validation"
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// PermissionsListQueryResolver for a list of permissions.
func (r *Resolver) PermissionsListQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := r.checkPermission(params.Context)
	if !allow {
		return []*domain.Permission{}, domain.ErrUnauthorized
	}

	permission, err := r.permissionUseCase.Fetch(params.Context)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return permission, nil
}

// PermissionQueryResolver for a single permission.
func (r *Resolver) PermissionQueryResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := r.checkPermission(params.Context)
	if !allow {
		return []*domain.Permission{}, domain.ErrUnauthorized
	}

	id, err := strconv.ParseInt(params.Args["ID"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	permission, err := r.permissionUseCase.GetByID(params.Context, id)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return permission, nil
}

// PermissionCreateResolver creates a new permission.
func (r *Resolver) PermissionCreateResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := r.checkPermission(params.Context)
	if !allow {
		return []*domain.Permission{}, domain.ErrUnauthorized
	}

	permission, err := storePermissionValidation(params)
	if err != nil {
		return nil, err
	}

	permission, err = r.permissionUseCase.Store(params.Context, permission)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return permission, nil
}

// PermissionUpdateResolver updates the given permission.
func (r *Resolver) PermissionUpdateResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := r.checkPermission(params.Context)
	if !allow {
		return []*domain.Permission{}, domain.ErrUnauthorized
	}

	permission, err := updatePermissionValidation(params)
	if err != nil {
		return nil, err
	}

	permission, err = r.permissionUseCase.Update(params.Context, permission)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return permission, nil
}

// PermissionDeleteResolver deletes the given permission.
func (r *Resolver) PermissionDeleteResolver(params graphql.ResolveParams) (interface{}, error) {
	allow := r.checkPermission(params.Context)
	if !allow {
		return []*domain.Permission{}, domain.ErrUnauthorized
	}

	id, err := strconv.ParseInt(params.Args["ID"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	err = r.permissionUseCase.Delete(params.Context, id)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}
	return nil, nil
}

func storePermissionValidation(params graphql.ResolveParams) (*domain.Permission, error) {
	permissionParams := params.Args["Permission"].(map[string]interface{})

	permission := &domain.Permission{
		Name:        permissionParams["name"].(string),
		Description: permissionParams["description"].(string),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := validation.IsAValidSchema(params.Context, permission); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return permission, nil
}

func updatePermissionValidation(params graphql.ResolveParams) (*domain.Permission, error) {
	permissionParams := params.Args["Permission"].(map[string]interface{})

	id, err := strconv.ParseInt(permissionParams["id"].(string), 10, 64)
	if err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	permission := &domain.Permission{
		ID:          id,
		Name:        permissionParams["name"].(string),
		Description: permissionParams["description"].(string),
		UpdatedAt:   time.Now(),
	}

	if err := validation.IsAValidSchema(params.Context, permission); err != nil {
		log.Error().Stack().Msg(err.Error())
		return nil, err
	}

	return permission, nil
}
