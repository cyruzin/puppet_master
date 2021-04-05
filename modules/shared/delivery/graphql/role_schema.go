package gql

import (
	"github.com/graphql-go/graphql"
)

var roleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Role",
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

var roleInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "RoleInput",
	Description: "Role payload for creating a new role",
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

var assingRoleToUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AssignRoleToUser",
	Fields: graphql.Fields{
		"role_id": &graphql.Field{
			Type: graphql.Int,
		},
		"user_id": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var assingRoleToUserTypeInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "AssingRoleToUserTypeInput",
	Description: "Assign/Sync role to a user",
	Fields: graphql.InputObjectConfigFieldMap{
		"role_id": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"user_id": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
})
