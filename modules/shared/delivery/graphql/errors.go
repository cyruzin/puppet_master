package gql

import "github.com/cyruzin/puppet_master/domain"

var (
	// ErrPermission will throw if the user does not have sufficient permission.
	ErrPermission = domain.ErrUnauthorized
)
