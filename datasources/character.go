package datasources

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Character struct {
	ID     int32 `gorm:"column:ID"`
	Name   string
	RaceID int32 `gorm:"column:raceID"`
	Race   Race  `gorm:"foreignKey:raceID"`
	Str    int32
	Dex    int32
	Con    int32
	Int    int32
	Wis    int32
	Cha    int32
	HP     int32 `gorm:"column:HP"`
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

func (r *Resolver) Character(ctx context.Context, args struct{ ID int32 }) *CharacterResolver {
	var character Character
	result := r.DB.Preload("Race").First(&character, args.ID)
	if result.Error != nil {
		return nil
	}

	return &CharacterResolver{c: &character}
}

func (r *Resolver) Characters() *[]*CharacterResolver {
	var characters []*Character

	result := r.DB.Preload("Race").Find(&characters)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	var xCharacterResolver []*CharacterResolver
	for _, c := range characters {
		xCharacterResolver = append(xCharacterResolver, &CharacterResolver{c})
	}

	return &xCharacterResolver
}
