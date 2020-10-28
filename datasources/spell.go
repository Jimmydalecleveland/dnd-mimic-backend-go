package datasources

import (
	"context"
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

func (r *SpellResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.s.ID)
}

func (r *SpellResolver) Name() string {
	return r.s.Name
}

func (r *Resolver) Spell(ctx context.Context, args struct{ ID int32 }) *SpellResolver {
	spellQuery := `
		Select "ID", name FROM "Spell"
		WHERE "ID" = $1
	`
	var spell Spell
	err := r.DB.QueryRow(spellQuery, args.ID).Scan(&spell.ID, &spell.Name)
	if err != nil {
		return nil
	}
	return &SpellResolver{s: &spell}
}

func (r *Resolver) Spells() *[]*SpellResolver {
	var spells []*Spell
	var err error
	spellQuery := `
		Select "ID", name FROM "Spell"
	`
	rows, err := r.DB.Query(spellQuery)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var singleSpell Spell
		err = rows.Scan(&singleSpell.ID, &singleSpell.Name)
		if err != nil {
			log.Fatal(err)
		}
		spells = append(spells, &singleSpell)
	}

	// forgive me father, I know not what else to do
	var xSpellResolver []*SpellResolver
	for _, s := range spells {
		xSpellResolver = append(xSpellResolver, &SpellResolver{s})
	}

	return &xSpellResolver
}
