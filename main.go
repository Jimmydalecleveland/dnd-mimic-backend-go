package main

import (
	"context"
	"log"
	"net/http"

	"github.com/friendsofgo/graphiql"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

const Schema = `
type Query {
	hello: String!
	spell(ID: ID!): Spell
}

type Spell {
	ID: ID!
	name: String!
}
`

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

func init() {
	for _, s := range spells {
		spellData[s.ID] = s
	}
}

func main() {
	schema := graphql.MustParseSchema(Schema, &query{})
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/query")
	if err != nil {
		panic(err)
	}

	http.Handle("/graphiql", graphiqlHandler)
	http.Handle("/query", &relay.Handler{Schema: schema})
	log.Println("Server up at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
