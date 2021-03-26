package gql

import "errors"

var (
	// ErrPermission will throw if the user does not have sufficient permission.
	ErrPermission = errors.New("insufficient permission")
)
