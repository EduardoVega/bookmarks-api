package gql

import "github.com/graphql-go/graphql"

// var tagType = graphql.NewObject(
//     graphql.ObjectConfig{
//         Name: "Tag",
//         Fields: graphql.Fields{
//             "name": &graphql.Field{
//                 Type: graphql.String,
//             },
//         },
//     },
// )

var bookmarkType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Bookmark",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.ID,
            },
            "name": &graphql.Field{
                Type: graphql.String,
            },
            "url": &graphql.Field{
                Type: graphql.String,
            },
            "tags": &graphql.Field{
                Type: graphql.NewList(graphql.String),
            },
        },
    },
)