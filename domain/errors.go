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
)
