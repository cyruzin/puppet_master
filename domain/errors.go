package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal server error")
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
	// ErrNotFound will throw if the requested resource could not be found
	ErrNotFound = errors.New("the resource you requested could not be found")
	// ErrAssignRole will throw if failed to assign role
	ErrAssignRole = errors.New("failed to assign role")
	// ErrRemoveRole will throw if failed to remove role
	ErrRemoveRole = errors.New("failed to remove role")
	// ErrSyncRole will throw if failed to sync role
	ErrSyncRole = errors.New("failed to sync role")
	// ErrPermissionByID will throw if failed to fetch permissions
	ErrPermissionByID = errors.New("failed to fetch permissions")
	// ErrAssignPermission will throw if failed to assign permission
	ErrAssignPermission = errors.New("failed to assign permission")
	// ErrRemovePermission will throw if failed to remove permission
	ErrRemovePermission = errors.New("failed to remove permission")
	// ErrSyncPermission will throw if failed to sync permission
	ErrSyncPermission = errors.New("failed to sync permission")
)
