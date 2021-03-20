package gql

import (
	"context"
	"strconv"

	"github.com/cyruzin/puppet_master/domain"
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// Resolver ...
type Resolver struct {
	userUseCase domain.UserUsecase
}

// UserResolver ...
func (r *Resolver) UserResolver(params graphql.ResolveParams) (interface{}, error) {
	idQuery, isOK := params.Args["id"].(string)

	parsedID, _ := strconv.ParseInt(idQuery, 10, 64)

	if isOK {
		user, err := r.userUseCase.GetByID(context.Background(), parsedID)
		if err != nil {
			log.Error().Stack().Err(err)
			return nil, err
		}
		return user, nil
	}
	return &domain.User{}, nil
}
