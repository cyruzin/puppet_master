package gql

import (
	"context"

	"github.com/cyruzin/puppet_master/domain"
)

func (r *Resolver) checkPermission(ctx context.Context) bool {
	auth := ctx.Value(domain.ContextKeyID).(*domain.Auth)

	if len(auth.UserPermissions) == 0 {
		return false
	}

	userPermissions, err := r.permissionUseCase.GetPermissionsByUserID(ctx, auth.UserID)
	if err != nil {
		return false
	}

	for _, role := range auth.Roles {
		if role == "admin" {
			return true
		}
	}

	for index, userPermission := range userPermissions {
		if auth.UserPermissions[index] != userPermission.Name {
			return false
		}
	}

	return true
}
