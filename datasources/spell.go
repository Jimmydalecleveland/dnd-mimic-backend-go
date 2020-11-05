package datasources

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
	_ "github.com/lib/pq"
)

type Spell struct {
	ID   int32 `gorm:"column:ID"`
	Name string
}

func (Spell) TableName() string {
	return "Spell"
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
	var spell Spell
	result := r.DB.First(&spell, args.ID)
	if result.Error != nil {
		return nil
	}
	return &SpellResolver{s: &spell}
}

func (r *Resolver) Spells() *[]*SpellResolver {
	var spells []*Spell

	result := r.DB.Find(&spells)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	// forgive me father, I know not what else to do
	var xSpellResolver []*SpellResolver
	for _, s := range spells {
		xSpellResolver = append(xSpellResolver, &SpellResolver{s})
	}

	return &xSpellResolver
}
