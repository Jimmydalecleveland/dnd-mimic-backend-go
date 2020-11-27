package datasources

import (
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Weapon struct {
	ItemID    int32
	Damage    *string
	SkillType string
	RangeType string
}

type WeaponResolver struct {
	w *Weapon
}

func (r *WeaponResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.w.ItemID)
}

func (r *WeaponResolver) Damage() *string {
	return r.w.Damage
}
func (r *WeaponResolver) SkillType() string {
	return r.w.SkillType
}
func (r *WeaponResolver) RangeType() string {
	return r.w.RangeType
}

func (r *Resolver) Weapons() *[]*WeaponResolver {
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
		err = rows.Scan(
			&weapon.ItemID,
			&weapon.Damage,
			&weapon.SkillType,
			&weapon.RangeType,
		)
		if err != nil {
			log.Fatal(err)
		}
		weapons = append(weapons, &weapon)
	}

	var xWeaponResolver []*WeaponResolver
	for _, w := range weapons {
		xWeaponResolver = append(xWeaponResolver, &WeaponResolver{w})
	}

	return &xWeaponResolver
}
