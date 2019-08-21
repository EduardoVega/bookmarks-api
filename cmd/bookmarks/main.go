package main

import (
	"log"
	"context"
	"net/http"
	"bookmarks/pkg/gql"
	gqlhandler "github.com/graphql-go/graphql-go-handler"
	"bookmarks/pkg/database"
)

func handler() (http.Handler, *database.MongoClient) {
	log.Println("=> Start bookmarks app")

	mongoClient, err := database.Connect("mongodb://bookmarks:Bookm%40rks@127.0.0.1:30543/bookmarks")

	// If mongo connection doesn't work, log error
	if err != nil{
		log.Fatal(err)
	}

	log.Println("=> Get GraphQL Schema")

	schema := gql.GetSchema(mongoClient)

	return gqlhandler.New(&gqlhandler.Config{
		Schema: schema.GQLSchema,
		Pretty: true,
	}), mongoClient

}

func main() {

	handler, mongoClient := handler()
	defer mongoClient.Disconnect(context.TODO())

	log.Println("=> Server listening")
	http.Handle("/bookmarks", handler)
	http.ListenAndServe(":8081", nil)

}