package gql

import (
    "github.com/graphql-go/graphql"
    "bookmarks/pkg/database"
)

// Query : struct that contains a graphql query object
type Query struct {
    QueryType *graphql.Object
}

// GetQuery : function that receives a mongo client object and returns a query struct
func GetQuery(m *database.MongoClient) *Query {

    resolver := Resolver{m}

    var queryType = graphql.NewObject(
        graphql.ObjectConfig{
            Name: "Query",
            Fields: graphql.Fields{
                "list": &graphql.Field{
                    Type: graphql.NewList(bookmarkType),
                    Resolve: resolver.GetBookmarksResolver,
                },
                "listByID": &graphql.Field{
                    Type: bookmarkType,
                    Resolve: resolver.GetBookmarkByIDResolver,
                    Args: graphql.FieldConfigArgument{
                        "id": &graphql.ArgumentConfig{
                            Type: graphql.NewNonNull(graphql.String),
                        },
                    },
                },
                "listByName": &graphql.Field{
                    Type: graphql.NewList(bookmarkType),
                    Resolve: resolver.GetBookmarksByNameResolver,
                    Args: graphql.FieldConfigArgument{
                        "name": &graphql.ArgumentConfig{
                            Type: graphql.NewNonNull(graphql.String),
                        },
                    },
                },
                "listByTags": &graphql.Field{
                    Type: graphql.NewList(bookmarkType),
                    Resolve: resolver.GetBookmarksByTagsResolver,
                    Args: graphql.FieldConfigArgument{
                        "tags": &graphql.ArgumentConfig{
                            Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
                        },
                    },
                },
            },
        },
    )

    return &Query{
        QueryType: queryType,
    }
}