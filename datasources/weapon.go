package datasources

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Weapon struct {
	ID        graphql.ID
	Name      string
	ItemType  string
	Damage    *string
	SkillType string
	RangeType string
	Cost      *string
	Weight    *string
}

type QuantifiedWeapon struct {
	Weapon
	Quantity int32
}

func (r *QuantifiedWeapon) ToQuantifiedWeapon() (*QuantifiedWeapon, bool) {
	return r, true
}

func (r *Resolver) Weapon(ctx context.Context, args struct{ ID int32 }) *Weapon {
	q := `
	SELECT 
		i."ID", 
		i.name, 
		i.type,
		w.damage, 
		w."skillType", 
		w."rangeType", 
		i.cost, 
		i.weight 
	FROM "Weapon" w
	JOIN "Item" i ON i."ID" = w."itemID"
	WHERE i."ID" = $1
 `
	var weapon Weapon
	var tempID int32
	err := r.DB.QueryRow(q, args.ID).Scan(
		&tempID,
		&weapon.Name,
		&weapon.ItemType,
		&weapon.Damage,
		&weapon.SkillType,
		&weapon.RangeType,
		&weapon.Cost,
		&weapon.Weight,
	)
	if err != nil {
		log.Print("Error received during Weapon query -> ", err)
		return nil
	}
	weapon.ID = Int32ToGraphqlID(tempID)

	return &weapon
}

func (r *Resolver) Weapons() *[]*Weapon {
	q := `
		SELECT 
			i."ID", 
			i.name, 
			i.type,
			w.damage, 
			w."skillType", 
			w."rangeType", 
			i.cost, 
			i.weight 
		FROM "Weapon" w
		JOIN "Item" i ON i."ID" = w."itemID"
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
			&weapon.Name,
			&weapon.ItemType,
			&weapon.Damage,
			&weapon.SkillType,
			&weapon.RangeType,
			&weapon.Cost,
			&weapon.Weight,
		)
		if err != nil {
			log.Fatal(err)
		}
		weapon.ID = Int32ToGraphqlID(tempID)
		weapons = append(weapons, &weapon)
	}

	return &weapons
}
