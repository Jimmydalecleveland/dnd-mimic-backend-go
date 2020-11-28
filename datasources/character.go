package datasources

import (
	"context"
	"database/sql"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Character struct {
	ID           int32
	Name         string
	MaxHP        int32
	HP           int32
	Str          int32
	Dex          int32
	Con          int32
	Int          int32
	Wis          int32
	Cha          int32
	Gp           int32
	Sp           int32
	Cp           int32
	Ep           int32
	Pp           int32
	RaceID       int32
	SubraceID    *int32
	BackgroundID int32
	CharClassID  int32
	DeathSaves   DeathSaves
	// UserID       int32
	// SpecID       int32
}

type DeathSaves struct {
	Successes int32
	Failures  int32
}

type CharacterResolver struct {
	character Character
	db        *sql.DB
}

type DeathSavesResolver struct {
	d DeathSaves
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

func (r *CharacterResolver) DeathSaves() *DeathSavesResolver {
	return &DeathSavesResolver{d: r.character.DeathSaves}
}

func (r *DeathSavesResolver) Successes() int32 {
	return r.d.Successes
}

func (r *DeathSavesResolver) Failures() int32 {
	return r.d.Failures
}

func (r *CharacterResolver) Race(ctx context.Context) *RaceResolver {
	race, err := queryRace(r.db, r.character.RaceID)
	if err != nil {
		return nil
	}
	subraces := querySubraces(r.db, r.character.RaceID)
	return &RaceResolver{r: race, s: subraces}
}

func (r *CharacterResolver) Subrace(ctx context.Context) *RaceResolver {
	if r.character.SubraceID == nil {
		return nil
	}
	subrace, err := querySubrace(r.db, *r.character.SubraceID)
	if err != nil {
		return nil
	}
	return &RaceResolver{r: subrace}
}

func (r *CharacterResolver) Background(ctx context.Context) *BackgroundResolver {
	b := queryBackground(r.db, r.character.BackgroundID)
	return &BackgroundResolver{b}
}

func (r *CharacterResolver) Class(ctx context.Context) *ClassResolver {
	c := queryClass(r.db, r.character.CharClassID)
	return &ClassResolver{c}
}

func (r *CharacterResolver) Skills() *[]*SkillResolver {
	var skills []*Skill

	charSkillQuery := `
	SELECT "Skill".* FROM "CharSkillProficiency"
	JOIN "Skill" ON "Skill"."ID" = "CharSkillProficiency"."skillID"
	WHERE "CharSkillProficiency"."charID" = $1
	`

	rows, err := r.db.Query(charSkillQuery, r.character.ID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var skill Skill
		err = rows.Scan(&skill.ID, &skill.Name, &skill.Ability)
		if err != nil {
			log.Fatal(err)
		}
		skills = append(skills, &skill)
	}

	var xSkillResolver []*SkillResolver
	for _, s := range skills {
		xSkillResolver = append(xSkillResolver, &SkillResolver{s})
	}
	return &xSkillResolver
}

func (r *CharacterResolver) Inventory() *[]*QuantifiedWeapon {
	q := `
		SELECT i."ID", i.name, w.damage, w."skillType", w."rangeType", i.cost, i.weight, ci.quantity FROM "Weapon" w
		JOIN "Item" i ON i."ID" = w."itemID"
		JOIN "CharacterItem" ci ON ci."itemID" = i."ID"
		WHERE ci."characterID" = $1
	`

	rows, err := r.db.Query(q, r.character.ID)
	if err != nil {
		log.Fatal(err)
	}

	var weapons []*QuantifiedWeapon
	for rows.Next() {
		var weapon QuantifiedWeapon
		var tempID int32
		err = rows.Scan(
			&tempID,
			&weapon.Name,
			&weapon.Damage,
			&weapon.SkillType,
			&weapon.RangeType,
			&weapon.Cost,
			&weapon.Weight,
			&weapon.Quantity,
		)
		if err != nil {
			log.Fatal(err)
		}
		weapon.ID = Int32ToGraphqlID(tempID)
		weapons = append(weapons, &weapon)
	}

	return &weapons
}

func (r *Resolver) Character(ctx context.Context, args struct{ ID int32 }) *CharacterResolver {
	q := `
	Select 
		c."ID",
		c.name, 
		c."maxHP", 
		c."HP", 
		c.str, 
		c.dex, 
		c.con, 
		c.int, 
		c.wis, 
		c.cha, 
		c.gp, 
		c.sp, 
		c.cp, 
		c.ep, 
		c.pp, 
		c."deathSaveSuccesses",
		c."deathSaveFailures",
		c."raceID", 
		c."subraceID", 
		c."backgroundID", 
		c."charClassID"
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
			&character.DeathSaves.Successes,
			&character.DeathSaves.Failures,
			&character.RaceID,
			&character.SubraceID,
			&character.BackgroundID,
			&character.CharClassID,
		)

	if err != nil {
		// TODO: figure out how to return empty object to GraphQL, or some non-breaking behavior
		// when no record is found
		log.Fatal(err)
	}

	return &CharacterResolver{character: character, db: r.DB}
}

func (r *Resolver) Characters() *[]*CharacterResolver {
	q := `
	Select 
		c."ID",
		c.name, 
		c."maxHP", 
		c."HP", 
		c.str, 
		c.dex, 
		c.con, 
		c.int, 
		c.wis, 
		c.cha, 
		c.gp, 
		c.sp, 
		c.cp, 
		c.ep, 
		c.pp, 
		c."deathSaveSuccesses",
		c."deathSaveFailures",
		c."raceID", 
		c."subraceID", 
		c."backgroundID", 
		c."charClassID"
	FROM "Character" c
	`
	var characters []Character

	rows, err := r.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var character Character
		err = rows.
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
				&character.DeathSaves.Successes,
				&character.DeathSaves.Failures,
				&character.RaceID,
				&character.SubraceID,
				&character.BackgroundID,
				&character.CharClassID,
			)

		if err != nil {
			log.Fatal(err)
		}
		characters = append(characters, character)
	}

	var xCharacterResolver []*CharacterResolver
	for _, character := range characters {
		xCharacterResolver = append(xCharacterResolver, &CharacterResolver{character: character, db: r.DB})
	}

	return &xCharacterResolver
}
