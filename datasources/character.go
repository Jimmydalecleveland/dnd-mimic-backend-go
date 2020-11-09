package datasources

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Character struct {
	ID    int32
	Name  string
	MaxHP int32
	HP    int32
	Str   int32
	Dex   int32
	Con   int32
	Int   int32
	Wis   int32
	Cha   int32
	Gp    int32
	Sp    int32
	Cp    int32
	Ep    int32
	Pp    int32
	// CharClassID  int32
	// RaceID       int32
	// UserID       int32
	// BackgroundID int32
	// SpecID       int32
	// SubraceID    int32
	// Deathsaves   string
	// Race   Race
	// Skills []Skill
}

func (Character) TableName() string {
	return "Character"
}

type CharacterResolver struct {
	c       *Character
	race    *Race
	subrace *Race
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

func (r *CharacterResolver) MaxHP() *int32 {
	return &r.c.MaxHP
}

func (r *CharacterResolver) HP() *int32 {
	return &r.c.HP
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
	return &RaceResolver{r: r.race}
}

func (r *CharacterResolver) Subrace() *RaceResolver {
	return &RaceResolver{r: r.subrace}
}

// func (r *CharacterResolver) Skills() *[]*SkillResolver {
// 	var xSkillResolver []*SkillResolver
// 	for _, s := range r.c.Skills {
// 		xSkillResolver = append(xSkillResolver, &SkillResolver{&s})
// 	}
// 	return &xSkillResolver
// }

func (r *Resolver) Character(ctx context.Context, args struct{ ID int32 }) *CharacterResolver {
	q := `
		SELECT c."ID", c.name, c."maxHP", c."HP", c.str, c.dex, c.con, c.int, c.wis, c.cha, c.gp, c.sp, c.cp, c.ep, c.pp, c."raceID", c."subraceID"
		FROM "Character" c
		WHERE "ID" = $1
	`

	var character Character
	var raceID int32
	var subraceID int32

	err := r.DB.
		QueryRow(q, args.ID).
		Scan(
			&character.ID,
			&character.Name,
			&character.MaxHP,
			&character.HP,
			&character.Str,
			&character.Dex,
			&character.Con,
			&character.Int,
			&character.Wis,
			&character.Cha,
			&character.Gp,
			&character.Sp,
			&character.Cp,
			&character.Ep,
			&character.Pp,
			&raceID,
			&subraceID,
		)

	if err != nil {
		// TODO: figure out how to return empty object to GraphQL, or some non-breaking behavior
		// when no record is found
		log.Fatal(err)
	}

	q2 := `Select "Race"."ID", "Race".name FROM "Race" WHERE "ID" = $1`
	var race Race
	err = r.DB.QueryRow(q2, raceID).Scan(&race.ID, &race.Name)
	if err != nil {
		log.Fatal(err)
	}

	var subrace Race
	err = r.DB.QueryRow(q2, subraceID).Scan(&subrace.ID, &subrace.Name)
	if err != nil {
		log.Fatal(err)
	}

	return &CharacterResolver{c: &character, race: &race, subrace: &subrace}
}

// func (r *Resolver) Characters() *[]*CharacterResolver {
// 	q := `
// 	Select c."ID", c.name, c."maxHP", c."HP", c.str, c.dex, c.con, c.int, c.wis, c.cha, c.gp, c.sp, c.cp, c.ep, c.pp
// 	FROM "Character"
// 	`
// 	var characters []*Character

// 	rows, err := r.DB.Query(q)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for rows.Next() {
// 		var character Character
// 		err = rows.
// 			Scan(
// 				&character.ID,
// 				&character.Name,
// 				&character.MaxHP,
// 				&character.HP,
// 				&character.Str,
// 				&character.Dex,
// 				&character.Con,
// 				&character.Int,
// 				&character.Wis,
// 				&character.Cha,
// 				&character.Gp,
// 				&character.Sp,
// 				&character.Cp,
// 				&character.Ep,
// 				&character.Pp,
// 			)

// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		characters = append(characters, &character)
// 	}

// 	var xCharacterResolver []*CharacterResolver
// 	for _, c := range characters {
// 		xCharacterResolver = append(xCharacterResolver, &CharacterResolver{c})
// 	}

// 	return &xCharacterResolver
// }
