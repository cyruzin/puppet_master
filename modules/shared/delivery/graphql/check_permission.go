package gql

import (
	"context"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/rs/zerolog/log"
)

func checkPermission(ctx context.Context) bool {
	permissions := ctx.Value(domain.ContextKeyID).(string)

	log.Info().Msg(permissions)

	return true
}
