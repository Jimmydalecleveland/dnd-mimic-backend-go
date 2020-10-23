package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/friendsofgo/graphiql"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type query struct{}

func (_ *query) Hello() string { return "Hello, world!" }

type Spell struct {
	ID   graphql.ID
	Name string
}

type SpellResolver struct {
	s *Spell
}

var spells = []*Spell{
	{ID: "1", Name: "Fireball"},
	{ID: "2", Name: "Aid"},
}

var spellData = make(map[graphql.ID]*Spell)

func (r *SpellResolver) ID() graphql.ID {
	return r.s.ID
}

func (r *SpellResolver) Name() string {
	return r.s.Name
}

func (_ *query) Spell(ctx context.Context, args struct{ ID graphql.ID }) *SpellResolver {
	s, ok := spellData[args.ID]
	if ok {
		return &SpellResolver{s}
	}
	return nil
}

func (_ *query) Spells() *[]*SpellResolver {
	// forgive me father, I know not what else to do
	var xSpellResolver []*SpellResolver
	for key, _ := range spellData {
		xSpellResolver = append(xSpellResolver, &SpellResolver{s: spellData[key]})
	}

	return &xSpellResolver
}

func init() {
	for _, s := range spells {
		spellData[s.ID] = s
	}
}

func main() {
	// Read .graphql file for schema
	schemaToString, err := ioutil.ReadFile("schema/schema.graphql")
	if err != nil {
		log.Fatal(err)
	}

	// Setup GraphQL with schema cast to string
	schema := graphql.MustParseSchema(string(schemaToString), &query{})
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
