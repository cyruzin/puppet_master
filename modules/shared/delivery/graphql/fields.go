package gql

import (
	"github.com/graphql-go/graphql"
)

func (r *Resolver) queryFields() graphql.Fields {
	fields := graphql.Fields{
		"fetchUsers": &graphql.Field{
			Type:        graphql.NewList(userType),
			Description: "Get a list of users",
			Resolve:     r.UsersListQueryResolver,
		},
		"getUser": &graphql.Field{
			Type:        userType,
			Description: "Get a single user",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.UserQueryResolver,
		},
	}

	return fields
}

func (r *Resolver) mutationFields() graphql.Fields {
	fields := graphql.Fields{
		"createUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Type: userInput,
				},
			},
			Resolve: r.UserCreateResolver,
		},
		"updateUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Type: userInput,
				},
			},
			Resolve: r.UserUpdateResolver,
		},
	}

	return fields
}
