package datasources

import (
	"fmt"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type GearPack struct {
	ID     graphql.ID
	Name   string
	Cost   string
	Weight string
	Items  *[]*QuantifiedAdventuringGear
}

type GearPackItem struct {
	ID       graphql.ID
	Name     string
	ItemType string
	Cost     *string
	Weight   *string
	Quantity int32
}

func (r *Resolver) GearPacks() *[]*GearPack {
	q := `
	SELECT * FROM "GearPack"
	`

	rows, err := r.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	var gearPacks []*GearPack
	for rows.Next() {
		var gearPack GearPack
		var tempID int32
		err = rows.Scan(
			&tempID,
			&gearPack.Name,
			&gearPack.Cost,
			&gearPack.Weight,
		)
		if err != nil {
			log.Fatal(err)
		}
		gearPack.ID = Int32ToGraphqlID(tempID)
		gearPacks = append(gearPacks, &gearPack)
	}

	q2 := `
		SELECT "Item".*, "GearPackItem".quantity FROM "GearPack"
		INNER JOIN "GearPackItem" ON "GearPackItem"."gearPackID" = "GearPack"."ID"
		INNER JOIN "Item" ON "Item"."ID" = "GearPackItem"."itemID"
		WHERE "GearPack"."ID" = $1
	`

	q3 := `
		SELECT 
			i."ID", 
			i.name, 
			i.type,
			a.description, 
			a.category, 
			a."categoryDescription", 
			i.cost, 
			i.weight 
		FROM "AdventuringGear" a
		JOIN "Item" i ON i."ID" = a."itemID"
		WHERE i."ID" = $1
 `

	for _, g := range gearPacks {
		rows, err = r.DB.Query(q2, g.ID)
		if err != nil {
			log.Fatal(err)
		}

		var finalGearPackItems []*QuantifiedAdventuringGear
		var gearPackItems []*GearPackItem
		for rows.Next() {
			var gearPackItem GearPackItem
			var tempID int32
			err = rows.Scan(
				&tempID,
				&gearPackItem.Name,
				&gearPackItem.ItemType,
				&gearPackItem.Cost,
				&gearPackItem.Weight,
				&gearPackItem.Quantity,
			)
			if err != nil {
				log.Fatal(err)
			}
			gearPackItem.ID = Int32ToGraphqlID(tempID)
			gearPackItems = append(gearPackItems, &gearPackItem)
		}

		for _, gi := range gearPackItems {
			if gi.ItemType == "AdventuringGear" {
				fmt.Println("Got an adventuring gear", gi.ID)
				var adventuringGear QuantifiedAdventuringGear
				var tempID int32
				err = r.DB.QueryRow(q3, gi.ID).Scan(
					&tempID,
					&adventuringGear.Name,
					&adventuringGear.ItemType,
					&adventuringGear.Description,
					&adventuringGear.Category,
					&adventuringGear.CategoryDescription,
					&adventuringGear.Cost,
					&adventuringGear.Weight,
				)
				adventuringGear.ID = Int32ToGraphqlID(tempID)
				adventuringGear.Quantity = gi.Quantity
				fmt.Println(adventuringGear)
				finalGearPackItems = append(finalGearPackItems, &adventuringGear)
			}
		}

		g.Items = &finalGearPackItems
	}

	return &gearPacks
}
