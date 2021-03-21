package gql

import (
	"github.com/graphql-go/graphql"
)

func (r *Resolver) queryFields() graphql.Fields {
	fields := graphql.Fields{
		// User
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

		// Role
		"fetchRoles": &graphql.Field{
			Type:        graphql.NewList(roleType),
			Description: "Get a list of roles",
			Resolve:     r.RolesListQueryResolver,
		},
		"getRole": &graphql.Field{
			Type:        roleType,
			Description: "Get a single role",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.RoleQueryResolver,
		},
	}

	return fields
}

func (r *Resolver) mutationFields() graphql.Fields {
	fields := graphql.Fields{
		// User
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
		"deleteUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.UserDeleteResolver,
		},

		// Role
		"createRole": &graphql.Field{
			Type: roleType,
			Args: graphql.FieldConfigArgument{
				"role": &graphql.ArgumentConfig{
					Type: roleInput,
				},
			},
			Resolve: r.RoleCreateResolver,
		},
		"updateRole": &graphql.Field{
			Type: roleType,
			Args: graphql.FieldConfigArgument{
				"role": &graphql.ArgumentConfig{
					Type: roleInput,
				},
			},
			Resolve: r.RoleUpdateResolver,
		},
		"deleteRole": &graphql.Field{
			Type: roleType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.RoleDeleteResolver,
		},
	}

	return fields
}
