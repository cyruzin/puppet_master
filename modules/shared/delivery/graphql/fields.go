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
		"RefreshToken": &graphql.Field{
			Type:        authType,
			Description: "Refreshes the token",
			Args: graphql.FieldConfigArgument{
				"UserID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: r.AuthRefreshTokenResolver,
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
		"GetPermissionsByRoleID": &graphql.Field{
			Type:        graphql.NewList(permissionType),
			Description: "Get all permissions by role ID",
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: r.PermissionGetByRoleIDResolver,
		},
		"GetPermissionsByRoleName": &graphql.Field{
			Type:        graphql.NewList(permissionType),
			Description: "Get all permissions by role name",
			Args: graphql.FieldConfigArgument{
				"Name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.PermissionGetByRoleNameResolver,
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
		"GetRolesByUserID": &graphql.Field{
			Type:        roleType,
			Description: "Get role by user ID",
			Args: graphql.FieldConfigArgument{
				"ID": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: r.RoleGetByUserIDResolver,
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
			Type: permissionRoleType,
			Args: graphql.FieldConfigArgument{
				"Permission": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(permissionRoleInput),
				},
			},
			Resolve: r.PermissionGiveResolver,
		},
		"SyncPermissionToRole": &graphql.Field{
			Type: permissionRoleType,
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
		"AssignRoleToUser": &graphql.Field{
			Type: assingRoleToUserType,
			Args: graphql.FieldConfigArgument{
				"Role": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(assingRoleToUserTypeInput),
				},
			},
			Resolve: r.RoleAssignResolver,
		},
		"SyncRoleToUser": &graphql.Field{
			Type: assingRoleToUserType,
			Args: graphql.FieldConfigArgument{
				"Role": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(assingRoleToUserTypeInput),
				},
			},
			Resolve: r.RoleSyncResolver,
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
