package gql

import (
	"github.com/cyruzin/puppet_master/domain"
	"github.com/graphql-go/graphql"
)

// Root holds a pointer to a graphql object
type Root struct {
	Query    *graphql.Object
	Mutation *graphql.Object
}

// Resolver struct for all use cases.
type Resolver struct {
	permissionUseCase domain.PermissionUsecase
	roleUseCase       domain.RoleUsecase
	userUseCase       domain.UserUsecase
}

func NewRoot(
	p domain.PermissionUsecase,
	r domain.RoleUsecase,
	u domain.UserUsecase,
) *Root {
	resolver := Resolver{
		permissionUseCase: p,
		roleUseCase:       r,
		userUseCase:       u,
	}
	root := Root{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: resolver.queryFields(),
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: resolver.mutationFields(),
		}),
	}

	return &root
}
