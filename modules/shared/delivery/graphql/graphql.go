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
	authUseCase       domain.AuthUsecase
	permissionUseCase domain.PermissionUsecase
	roleUseCase       domain.RoleUsecase
	userUseCase       domain.UserUsecase
}

func NewRoot(
	auth domain.AuthUsecase,
	permission domain.PermissionUsecase,
	role domain.RoleUsecase,
	user domain.UserUsecase,
) *Root {
	resolver := Resolver{
		authUseCase:       auth,
		permissionUseCase: permission,
		roleUseCase:       role,
		userUseCase:       user,
	}
	root := Root{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:        "Queries",
			Fields:      resolver.queryFields(),
			Description: "All Puppet Master queries",
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:        "Mutations",
			Fields:      resolver.mutationFields(),
			Description: "All Puppet Master mutations",
		}),
	}

	return &root
}
