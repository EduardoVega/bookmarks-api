package gql

import (
    "github.com/graphql-go/graphql"
    "bookmarks/pkg/database"
)

// Mutation :
type Mutation struct {
	MutationType *graphql.Object
}

// GetMutation :
func GetMutation(m *database.MongoClient) *Mutation {

	resolver := Resolver{m}

	mutationType := graphql.NewObject(
		graphql.ObjectConfig {
			Name: "Mutation",
			Fields: graphql.Fields{
				"create": &graphql.Field{
					Type: bookmarkType,
					Description: "Create new bookmark",
					Args: graphql.FieldConfigArgument {
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"url": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"tags": &graphql.ArgumentConfig{
							Type: graphql.NewList(graphql.String),
						},
					},
					Resolve: resolver.CreateBookmarkResolver,
				},
				"update": &graphql.Field{
					Type: bookmarkType,
					Description: "Create new bookmark",
					Args: graphql.FieldConfigArgument {
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"url": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"tags": &graphql.ArgumentConfig{
							Type: graphql.NewList(graphql.String),
						},
					},
					Resolve: resolver.UpdateBookmarkResolver,
				},
				"delete": &graphql.Field{
					Type: graphql.String,
					Description: "Delete a bookmark",
					Args: graphql.FieldConfigArgument {
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: resolver.DeleteBookmarkResolver,
				},
			},
		},
	)

	return &Mutation{
		MutationType: mutationType,
	}
}