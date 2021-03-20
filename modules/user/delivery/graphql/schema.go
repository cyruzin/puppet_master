package gql

import (
	"github.com/cyruzin/puppet_master/domain"
	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"superadmin": &graphql.Field{
			Type: graphql.Boolean,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updated_at": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})

// Root holds a pointer to a graphql object
type Root struct {
	Query *graphql.Object
}

func NewRoot(u domain.UserUsecase) *Root {
	resolver := Resolver{userUseCase: u}

	root := Root{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{

				"user": &graphql.Field{
					Type:        userType,
					Description: "Get a single user",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},

					Resolve: resolver.UserResolver,
				},
			},
		}),
	}

	return &root
}
