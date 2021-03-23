package gql

import (
	"github.com/graphql-go/graphql"
)

func (r *Resolver) queryFields() graphql.Fields {
	fields := graphql.Fields{
		// Auth
		"authenticate": &graphql.Field{
			Type:        authType,
			Description: "Authenticate the given user",
			Args: graphql.FieldConfigArgument{
				"credentials": &graphql.ArgumentConfig{
					Type: authInput,
				},
			},
			Resolve: r.AuthQueryResolver,
		},

		// Permission
		"fetchPermissions": &graphql.Field{
			Type:        graphql.NewList(permissionType),
			Description: "Get a list of permissions",
			Resolve:     r.PermissionsListQueryResolver,
		},
		"getPermission": &graphql.Field{
			Type:        permissionType,
			Description: "Get a single permission",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.PermissionQueryResolver,
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
	}

	return fields
}

func (r *Resolver) mutationFields() graphql.Fields {
	fields := graphql.Fields{
		// Permission
		"createPermission": &graphql.Field{
			Type: permissionType,
			Args: graphql.FieldConfigArgument{
				"permission": &graphql.ArgumentConfig{
					Type: permissionInput,
				},
			},
			Resolve: r.PermissionCreateResolver,
		},
		"updatePermission": &graphql.Field{
			Type: permissionType,
			Args: graphql.FieldConfigArgument{
				"permission": &graphql.ArgumentConfig{
					Type: permissionInput,
				},
			},
			Resolve: r.PermissionUpdateResolver,
		},
		"deletePermission": &graphql.Field{
			Type: permissionType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.PermissionDeleteResolver,
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
	}

	return fields
}
