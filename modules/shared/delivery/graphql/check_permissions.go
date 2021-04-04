package gql

import (
	"context"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/rs/zerolog/log"
)

func (r *Resolver) checkPermissions(ctx context.Context, permission string, roles []string) bool {
	user := ctx.Value(domain.ContextKeyID).(map[string]interface{})

	userCache := &domain.UserCache{}

	if err := r.authUseCase.GetCache(ctx, user["email"].(string), userCache); err != nil {
		log.Error().Stack().Msg(err.Error())
		return false
	}

	if userCache.Role == "Admin" {
		return true
	}

	for _, currentRole := range roles {
		if currentRole == userCache.Role {
			return true
		}
	}

	for _, currentPermission := range userCache.Permissions {
		if currentPermission == permission {
			return true
		}
	}

	return false
}
