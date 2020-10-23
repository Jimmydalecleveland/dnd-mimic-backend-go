package datasources

import (
	"context"
	"database/sql"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
	_ "github.com/lib/pq"
)

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
	return Int32ToGraphqlID(r.s.ID)
}

func (r *SpellResolver) Name() string {
	return r.s.Name
}

func (_ *Query) Spell(ctx context.Context, args struct{ ID int32 }) *SpellResolver {
	s, ok := spellData[args.ID]
	if ok {
		return &SpellResolver{s}
	}
	return nil
}

func (_ *Query) Spells() *[]*SpellResolver {
	// forgive me father, I know not what else to do
	var xSpellResolver []*SpellResolver
	for key := range spellData {
		xSpellResolver = append(xSpellResolver, &SpellResolver{s: spellData[key]})
	}

	return &xSpellResolver
}

func QuerySpells(db *sql.DB) {
	spellQuery := `
		Select "ID", name FROM "Spell"
	`
	rows, rerr := db.Query(spellQuery)
	if rerr != nil {
		panic(rerr)
	}
	for rows.Next() {
		var singleSpell Spell
		err := rows.Scan(&singleSpell.ID, &singleSpell.Name)
		if err != nil {
			log.Fatal(err)
		}
		spells = append(spells, &singleSpell)
	}

	for _, s := range spells {
		spellData[s.ID] = s
	}
}
