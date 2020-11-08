package datasources

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Character struct {
	ID     int32
	Name   string
	RaceID int32
	Race   Race
	Skills []Skill
	Str    int32
	Dex    int32
	Con    int32
	Int    int32
	Wis    int32
	Cha    int32
	HP     int32
	Gp     int32
	Sp     int32
	Cp     int32
}

func (Character) TableName() string {
	return "Character"
}

type CharacterResolver struct {
	c *Character
}

func (r *CharacterResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.c.ID)
}

func (r *CharacterResolver) Name() *string {
	name := r.c.Name
	// should we return nil for gql?
	if name == "" {
		return nil
	}
	return &name
}

func (r *CharacterResolver) Str() *int32 {
	return &r.c.Str
}

func (r *CharacterResolver) Dex() *int32 {
	return &r.c.Dex
}

func (r *CharacterResolver) Con() *int32 {
	return &r.c.Con
}

func (r *CharacterResolver) Int() *int32 {
	return &r.c.Int
}

func (r *CharacterResolver) Wis() *int32 {
	return &r.c.Wis
}

func (r *CharacterResolver) Cha() *int32 {
	return &r.c.Cha
}

func (r *CharacterResolver) HP() *int32 {
	return &r.c.HP
}

func (r *CharacterResolver) Gp() *int32 {
	return &r.c.Gp
}

func (r *CharacterResolver) Sp() *int32 {
	return &r.c.Sp
}

func (r *CharacterResolver) Cp() *int32 {
	return &r.c.Cp
}

func (r *CharacterResolver) Race() *RaceResolver {
	return &RaceResolver{r: &r.c.Race}
}

func (r *CharacterResolver) Skills() *[]*SkillResolver {
	var xSkillResolver []*SkillResolver
	for _, s := range r.c.Skills {
		xSkillResolver = append(xSkillResolver, &SkillResolver{&s})
	}
	return &xSkillResolver
}

func (r *Resolver) Character(ctx context.Context, args struct{ ID int32 }) *CharacterResolver {
	q := `Select * FROM "Character" WHERE "ID" = $1`
	var character Character
	err := r.DB.QueryRow(q, args.ID).Scan(&character.ID, &character.Name)
	if err != nil {
		log.Fatal(err)
	}

	return &CharacterResolver{c: &character}
}

func (r *Resolver) Characters() *[]*CharacterResolver {
	q := `SELECT * FROM "Character"`
	var characters []*Character

	rows, err := r.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var character Character
		err = rows.Scan(&character.ID, &character.Name)
		if err != nil {
			log.Fatal(err)
		}
		characters = append(characters, &character)
	}

	var xCharacterResolver []*CharacterResolver
	for _, c := range characters {
		xCharacterResolver = append(xCharacterResolver, &CharacterResolver{c})
	}

	return &xCharacterResolver
}
