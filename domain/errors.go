package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal server error")
	// ErrNotFound will throw if the requested resource could not be found
	ErrNotFound = errors.New("the resource you requested could not be found")
	// ErrUnauthorized will throw if the user does not have sufficient permission
	ErrUnauthorized = errors.New("unauthorized")
	// ErrBadRequest will throw if the user send a bad payload
	ErrBadRequest = errors.New("bad request")
	// ErrIDParam will throw if the user do not provide a valid id
	ErrIDParam = errors.New("invalid id format")

	// ErrFetchError will throw if failed to fetch
	ErrFetchError = errors.New("failed to fetch")
	// ErrGetByIDError will throw if failed to get by id
	ErrGetByIDError = errors.New("failed to get by id")
	// ErrStoreError will throw if failed to store
	ErrStoreError = errors.New("failed to store")
	// ErrUpdateError will throw if failed to update
	ErrUpdateError = errors.New("failed to update")
	// ErrDeleteError will throw if failed to delete
	ErrDeleteError = errors.New("failed to delete")

	// ErrRoleByID will throw if failed to fetch roles by id
	ErrRoleByID = errors.New("failed to fetch roles by id")
	// ErrAssignRole will throw if failed to assign role
	ErrAssignRole = errors.New("failed to assign role")
	// ErrRemoveRole will throw if failed to remove role
	ErrRemoveRole = errors.New("failed to remove role")
	// ErrSyncRole will throw if failed to sync role
	ErrSyncRole = errors.New("failed to sync role")

	// ErrPermissionByID will throw if failed to fetch permissions by id
	ErrPermissionByID = errors.New("failed to fetch permissions by id")
	// ErrAssignPermission will throw if failed to assign permission
	ErrAssignPermission = errors.New("failed to assign permission")
	// ErrRemovePermission will throw if failed to remove permission
	ErrRemovePermission = errors.New("failed to remove permission")
	// ErrSyncPermission will throw if failed to sync permission
	ErrSyncPermission = errors.New("failed to sync permission")

	// ErrSetCache will throw if failed to set cache data
	ErrSetCache = errors.New("failed to set cache data")
	// ErrGetCache will throw if failed to get cache data
	ErrGetCache = errors.New("failed to get cache data")
	// ErrCacheKeyNil will throw if failed to find the cache key
	ErrCacheKeyNil = errors.New("failed to find the cache key")
	// ErrCacheMarshalling will throw if failed to marshal the cache
	ErrCacheMarshalling = errors.New("failed to marshal the cache")
	// ErrCacheUnmarshalling will throw if failed to unmarshal the cache key
	ErrCacheUnmarshalling = errors.New("failed to unmarshal the cache key")

	// ErrUserID will throw if the ID is invalid
	ErrUserID = errors.New("invalid user id")
)
