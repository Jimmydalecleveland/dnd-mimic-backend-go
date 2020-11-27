package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/friendsofgo/graphiql"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/jimmydalecleveland/dnd-mimic-backend-go/database"
	"github.com/jimmydalecleveland/dnd-mimic-backend-go/datasources"
)

func main() {
	// Open db connection
	db, err := database.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Read .graphql file for schema
	schemaToString, err := ioutil.ReadFile("./schema/schema.graphql")
	if err != nil {
		log.Fatal(err)
	}

	// Setup GraphQL with schema cast to string and db instance
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	schema := graphql.MustParseSchema(string(schemaToString), &datasources.Resolver{DB: db}, opts...)
	http.Handle("/query", &relay.Handler{Schema: schema})

	// Setup graphiql
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/query")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/graphiql", graphiqlHandler)

	log.Println("Server up at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
