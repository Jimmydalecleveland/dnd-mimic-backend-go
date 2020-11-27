package datasources

import (
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Weapon struct {
	ID        graphql.ID
	Damage    *string
	SkillType string
	RangeType string
}

func (r *Resolver) Weapons() *[]*Weapon {
	q := `
		SELECT * FROM "Weapon"
	`
	rows, err := r.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	var weapons []*Weapon
	for rows.Next() {
		var weapon Weapon
		var tempID int32
		err = rows.Scan(
			&tempID,
			&weapon.Damage,
			&weapon.SkillType,
			&weapon.RangeType,
		)
		if err != nil {
			log.Fatal(err)
		}
		weapon.ID = Int32ToGraphqlID(tempID)
		weapons = append(weapons, &weapon)
	}

	return &weapons
}
