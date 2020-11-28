package datasources

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Armor struct {
	ID                    graphql.ID
	Name                  string
	ItemType              string
	Ac                    int32
	IsDexAdded            bool
	DisadvantageOnStealth bool
	MaxDex                *int32
	Cost                  *string
	Weight                *string
}

type QuantifiedArmor struct {
	Armor
	Quantity int32
}

func (r *Resolver) Armor(ctx context.Context, args struct{ ID int32 }) *Armor {
	q := `
		SELECT 
			i."ID", 
			i.name, 
			i.type,
			a.ac, 
			a."isDexAdded", 
			a."disadvantageOnStealth", 
			a."maxDex",
			i.cost, 
			i.weight 
		FROM "Armor" a
		JOIN "Item" i ON i."ID" = a."itemID"
		WHERE i."ID" = $1
 `
	var armor Armor
	var tempID int32
	err := r.DB.QueryRow(q, args.ID).Scan(
		&tempID,
		&armor.Name,
		&armor.ItemType,
		&armor.Ac,
		&armor.IsDexAdded,
		&armor.DisadvantageOnStealth,
		&armor.MaxDex,
		&armor.Cost,
		&armor.Weight,
	)
	if err != nil {
		log.Print("Error received during Armor query -> ", err)
		return nil
	}
	armor.ID = Int32ToGraphqlID(tempID)

	return &armor
}

func (r *Resolver) Armors() *[]*Armor {
	q := `
		SELECT 
			i."ID", 
			i.name, 
			i.type,
			a.ac, 
			a."isDexAdded", 
			a."disadvantageOnStealth", 
			a."maxDex",
			i.cost, 
			i.weight 
		FROM "Armor" a
		JOIN "Item" i ON i."ID" = a."itemID"
	`
	rows, err := r.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	var armors []*Armor
	for rows.Next() {
		var armor Armor
		var tempID int32
		err = rows.Scan(
			&tempID,
			&armor.Name,
			&armor.ItemType,
			&armor.Ac,
			&armor.IsDexAdded,
			&armor.DisadvantageOnStealth,
			&armor.MaxDex,
			&armor.Cost,
			&armor.Weight,
		)
		if err != nil {
			log.Fatal(err)
		}
		armor.ID = Int32ToGraphqlID(tempID)
		armors = append(armors, &armor)
	}

	return &armors
}
