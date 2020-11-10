package datasources

import (
	"context"
	"database/sql"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Character struct {
	ID        int32
	Name      string
	MaxHP     int32
	HP        int32
	Str       int32
	Dex       int32
	Con       int32
	Int       int32
	Wis       int32
	Cha       int32
	Gp        int32
	Sp        int32
	Cp        int32
	Ep        int32
	Pp        int32
	RaceID    int32
	SubraceID int32
	// CharClassID  int32
	// UserID       int32
	// BackgroundID int32
	// SpecID       int32
	// Deathsaves   string
	// Race   Race
	// Skills []Skill
}

func (Character) TableName() string {
	return "Character"
}

type CharacterResolver struct {
	character Character
	db        *sql.DB
}

func (r *CharacterResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.character.ID)
}

func (r *CharacterResolver) Name() *string {
	name := r.character.Name
	// should we return nil for gql?
	if name == "" {
		return nil
	}
	return &name
}

func (r *CharacterResolver) MaxHP() *int32 {
	return &r.character.MaxHP
}

func (r *CharacterResolver) HP() *int32 {
	return &r.character.HP
}

func (r *CharacterResolver) Str() *int32 {
	return &r.character.Str
}

func (r *CharacterResolver) Dex() *int32 {
	return &r.character.Dex
}

func (r *CharacterResolver) Con() *int32 {
	return &r.character.Con
}

func (r *CharacterResolver) Int() *int32 {
	return &r.character.Int
}

func (r *CharacterResolver) Wis() *int32 {
	return &r.character.Wis
}

func (r *CharacterResolver) Cha() *int32 {
	return &r.character.Cha
}

func (r *CharacterResolver) Gp() *int32 {
	return &r.character.Gp
}

func (r *CharacterResolver) Sp() *int32 {
	return &r.character.Sp
}

func (r *CharacterResolver) Cp() *int32 {
	return &r.character.Cp
}

func (r *CharacterResolver) Race(ctx context.Context) *RaceResolver {
	race := queryRace(r.db, r.character.RaceID)
	subraces := querySubraces(r.db, r.character.RaceID)
	return &RaceResolver{r: race, s: subraces}
}

func (r *CharacterResolver) Subrace(ctx context.Context) *RaceResolver {
	if r.character.SubraceID == 0 {
		return nil
	}
	subrace := queryRace(r.db, r.character.SubraceID)
	return &RaceResolver{r: subrace}
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
			&character.RaceID,
			&character.SubraceID,
		)

	if err != nil {
		// TODO: figure out how to return empty object to GraphQL, or some non-breaking behavior
		// when no record is found
		log.Fatal(err)
	}

	return &CharacterResolver{character: character, db: r.DB}
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
