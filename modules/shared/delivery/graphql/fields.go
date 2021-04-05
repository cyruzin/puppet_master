package gql

import (
	"github.com/graphql-go/graphql"
)

func (r *Resolver) queryFields() graphql.Fields {
	fields := graphql.Fields{
		// Auth
		"Authenticate": &graphql.Field{
			Type:        authType,
			Description: "Authenticate the given user",
			Args: graphql.FieldConfigArgument{
				"Credentials": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(authInput),
				},
			},
			Resolve: r.AuthQueryResolver,
		},

		// Permission
		"FetchPermissions": &graphql.Field{
			Type:        graphql.NewList(permissionType),
			Description: "Get a list of permissions",
			Resolve:     r.PermissionsListQueryResolver,
		},
		"GetPermission": &graphql.Field{
			Type:        permissionType,
			Description: "Get a single permission",
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.PermissionQueryResolver,
		},

		// Role
		"FetchRoles": &graphql.Field{
			Type:        graphql.NewList(roleType),
			Description: "Get a list of roles",
			Resolve:     r.RolesListQueryResolver,
		},
		"GetRole": &graphql.Field{
			Type:        roleType,
			Description: "Get a single role",
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.RoleQueryResolver,
		},

		// User
		"FetchUsers": &graphql.Field{
			Type:        graphql.NewList(userType),
			Description: "Get a list of users",
			Resolve:     r.UsersListQueryResolver,
		},
		"GetUser": &graphql.Field{
			Type:        userType,
			Description: "Get a single user",
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
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
		"CreatePermission": &graphql.Field{
			Type: permissionType,
			Args: graphql.FieldConfigArgument{
				"Permission": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(permissionInput),
				},
			},
			Resolve: r.PermissionCreateResolver,
		},
		"UpdatePermission": &graphql.Field{
			Type: permissionType,
			Args: graphql.FieldConfigArgument{
				"Permission": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(permissionInput),
				},
			},
			Resolve: r.PermissionUpdateResolver,
		},
		"DeletePermission": &graphql.Field{
			Type: permissionType,
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.PermissionDeleteResolver,
		},
		"GivePermissionToRole": &graphql.Field{
			Type: permissionType,
			Args: graphql.FieldConfigArgument{
				"Permission": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(permissionRoleInput),
				},
			},
			Resolve: r.PermissionGiveResolver,
		},
		"RemovePermissionToRole": &graphql.Field{
			Type: permissionType,
			Args: graphql.FieldConfigArgument{
				"Permission": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(permissionRoleInput),
				},
			},
			Resolve: r.PermissionRemoveResolver,
		},
		"SyncPermissionToRole": &graphql.Field{
			Type: permissionType,
			Args: graphql.FieldConfigArgument{
				"Permission": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(permissionRoleInput),
				},
			},
			Resolve: r.PermissionSyncResolver,
		},

		// Role
		"CreateRole": &graphql.Field{
			Type: roleType,
			Args: graphql.FieldConfigArgument{
				"Role": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(roleInput),
				},
			},
			Resolve: r.RoleCreateResolver,
		},
		"UpdateRole": &graphql.Field{
			Type: roleType,
			Args: graphql.FieldConfigArgument{
				"Role": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(roleInput),
				},
			},
			Resolve: r.RoleUpdateResolver,
		},
		"DeleteRole": &graphql.Field{
			Type: roleType,
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.RoleDeleteResolver,
		},

		// User
		"CreateUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"User": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(userInput),
				},
			},
			Resolve: r.UserCreateResolver,
		},
		"UpdateUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"User": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(userInput),
				},
			},
			Resolve: r.UserUpdateResolver,
		},
		"DeleteUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.UserDeleteResolver,
		},
	}

	return fields
}
