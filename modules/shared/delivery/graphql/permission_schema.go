package gql

import (
	"github.com/graphql-go/graphql"
)

var permissionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Permission",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updated_at": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})

var permissionInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "PermissionInput",
	Description: "Permission payload for creating a new permission",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"created_at": &graphql.InputObjectFieldConfig{
			Type: graphql.DateTime,
		},
		"udpated_at": &graphql.InputObjectFieldConfig{
			Type: graphql.DateTime,
		},
	},
})

var permissionRoleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PermissionRole",
	Fields: graphql.Fields{
		"role_id": &graphql.Field{
			Type: graphql.Int,
		},
		"permissions": &graphql.Field{
			Type: &graphql.List{
				OfType: graphql.Int,
			},
		},
	},
})

var permissionRoleInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "PermissionRoleInput",
	Description: "Give/Sync permission to a role",
	Fields: graphql.InputObjectConfigFieldMap{
		"role_id": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"permissions": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(&graphql.List{
				OfType: graphql.Int,
			}),
		},
	},
})
