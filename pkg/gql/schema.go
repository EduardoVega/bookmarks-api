package gql

import (
    "github.com/graphql-go/graphql"
    "bookmarks/pkg/database"
)

// Schema : struct that contains a graphql schema
type Schema struct {
    GQLSchema *graphql.Schema
}

// GetSchema : function that returns a schema struct
func GetSchema(m *database.MongoClient) *Schema {

    query := GetQuery(m)
    mutation := GetMutation(m)

    var schema, _ = graphql.NewSchema(
        graphql.SchemaConfig{
            Query: query.QueryType,
            Mutation: mutation.MutationType,
    })

    return &Schema{&schema}
}