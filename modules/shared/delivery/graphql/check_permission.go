package gql

import (
	"context"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/rs/zerolog/log"
)

func checkPermission(ctx context.Context) bool {
	permissions := ctx.Value(domain.ContextKeyID).(string)

	if permissions == "" {
		return false
	}

	log.Info().Msg(permissions)

	return true
}
