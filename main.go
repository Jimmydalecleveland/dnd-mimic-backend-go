package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/friendsofgo/graphiql"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/jimmydalecleveland/go-graphql-server/datasources"
	"github.com/jimmydalecleveland/go-graphql-server/pgconnect"
)

func main() {
	// Open db connection
	db := pgconnect.InitializeDB()
	defer db.Close()

	// Read .graphql file for schema
	schemaToString, err := ioutil.ReadFile("schema/schema.graphql")
	if err != nil {
		log.Fatal(err)
	}

	// Setup GraphQL with schema cast to string
	schema := graphql.MustParseSchema(string(schemaToString), &datasources.Resolver{DB: db})
	http.Handle("/query", &relay.Handler{Schema: schema})

	// Setup graphiql
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/query")
	if err != nil {
		panic(err)
	}
	http.Handle("/graphiql", graphiqlHandler)

	log.Println("Server up at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
