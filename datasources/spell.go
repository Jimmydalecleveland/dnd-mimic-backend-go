package datasources

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Spell struct {
	ID          int32
	Name        string
	Level       int32
	School      string
	CastingTime string
	Range       string
	Components  string
	Duration    string
	Description string
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

func (r *SpellResolver) Level() *int32 {
	return &r.s.Level
}

func (r *SpellResolver) School() *string {
	return &r.s.School
}

func (r *SpellResolver) CastingTime() *string {
	return &r.s.CastingTime
}

func (r *SpellResolver) Range() *string {
	return &r.s.Range
}

func (r *SpellResolver) Components() *string {
	return &r.s.Components
}

func (r *SpellResolver) Duration() *string {
	return &r.s.Duration
}

func (r *SpellResolver) Description() *string {
	return &r.s.Description
}

func (r *Resolver) Spell(ctx context.Context, args struct{ ID int32 }) *SpellResolver {
	spellQuery := `
		Select * FROM "Spell"
		WHERE "ID" = $1
	`
	var spell Spell
	err := r.DB.
		QueryRow(spellQuery, args.ID).
		Scan(
			&spell.ID,
			&spell.Name,
			&spell.Level,
			&spell.School,
			&spell.CastingTime,
			&spell.Range,
			&spell.Components,
			&spell.Duration,
			&spell.Description,
		)

	if err != nil {
		return nil
	}
	return &SpellResolver{s: &spell}
}

func (r *Resolver) Spells() *[]*SpellResolver {
	var spells []*Spell

	var err error
	spellQuery := `
		Select * FROM "Spell"
	`
	rows, err := r.DB.Query(spellQuery)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var spell Spell
		err = rows.
			Scan(
				&spell.ID,
				&spell.Name,
				&spell.Level,
				&spell.School,
				&spell.CastingTime,
				&spell.Range,
				&spell.Components,
				&spell.Duration,
				&spell.Description,
			)

		if err != nil {
			log.Fatal(err)
		}
		spells = append(spells, &spell)
	}

	// forgive me father, I know not what else to do
	var xSpellResolver []*SpellResolver
	for _, s := range spells {
		xSpellResolver = append(xSpellResolver, &SpellResolver{s})
	}

	return &xSpellResolver
}
