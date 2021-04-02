package gql

import (
	"context"
	"fmt"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/rs/zerolog/log"
)

func (r *Resolver) checkPermissions(ctx context.Context, permission string) bool {
	user := ctx.Value(domain.ContextKeyID).(map[string]interface{})

	userCache := &domain.UserCache{}

	fmt.Println(user["email"].(string))

	if err := r.authUseCase.GetCache(ctx, user["email"].(string), userCache); err != nil {
		log.Error().Stack().Msg(err.Error())
		fmt.Println("erro do cache")
		return false
	}

	if userCache.Role == "Admin" {
		return true
	}

	for _, currentPermission := range userCache.Permissions {
		if currentPermission == permission {
			fmt.Println("permissão bate")
			return true
		}
	}

	fmt.Println("chegou no final")
	return false
}
