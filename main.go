package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/friendsofgo/graphiql"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// This is nutty, but converting from an int32 to a string requires it to first
// be converted to an int. This library expects a graphql.ID, which a string can
// be converted to. This method is benchmark faster than `fmt.Sprint(i)`
func int32ToGraphqlID(i int32) graphql.ID {
	return graphql.ID(strconv.Itoa(int(i)))
}

type query struct{}

func (_ *query) Hello() string { return "Hello, world!" }

type Spell struct {
	ID   int32
	Name string
}

type SpellResolver struct {
	s *Spell
}

var spells []*Spell

var spellData = make(map[int32]*Spell)

func (r *SpellResolver) ID() graphql.ID {
	return int32ToGraphqlID(r.s.ID)
}

func (r *SpellResolver) Name() string {
	return r.s.Name
}

func (_ *query) Spell(ctx context.Context, args struct{ ID int32 }) *SpellResolver {
	s, ok := spellData[args.ID]
	if ok {
		return &SpellResolver{s}
	}
	return nil
}

func (_ *query) Spells() *[]*SpellResolver {
	// forgive me father, I know not what else to do
	var xSpellResolver []*SpellResolver
	for key := range spellData {
		xSpellResolver = append(xSpellResolver, &SpellResolver{s: spellData[key]})
	}

	return &xSpellResolver
}

func db() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatalf("Error loading database .env file")
	}
	var (
		host     = os.Getenv("PGHOST")
		port, _  = strconv.Atoi(os.Getenv("PGPORT"))
		user     = os.Getenv("PGUSER")
		dbname   = os.Getenv("PGDATABASE")
		password = os.Getenv("PGPASSWORD")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to db")

	// sqlStatement := `
	// SELECT * FROM "Spell"
	// `
	// var dbSpells []Spell
	// rows, rerr := db.Query(sqlStatement)
	// if rerr != nil {
	// 	panic(err)
	// }

	spellQuery := `
		Select "ID", name FROM "Spell"
	`
	rows, rerr := db.Query(spellQuery)
	if rerr != nil {
		panic(rerr)
	}
	defer rows.Close()
	for rows.Next() {
		var singleSpell Spell
		err := rows.Scan(&singleSpell.ID, &singleSpell.Name)
		if err != nil {
			log.Fatal(err)
		}
		spells = append(spells, &singleSpell)
	}
}

func init() {
	db()
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
