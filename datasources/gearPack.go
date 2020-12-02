package datasources

import (
	"fmt"
	"log"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
)

type GearPack struct {
	ID     graphql.ID
	Name   string
	Cost   string
	Weight string
	Items  *[]*ItemResolver
}

type QuantifiedTool struct {
	ID          graphql.ID
	Name        string
	ItemType    string
	Category    *string
	Description *string
	Cost        *string
	Weight      *string
	Quantity    int32
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

	gearPackItemQuery := `
		SELECT "Item"."ID", "Item".type FROM "GearPackItem"
		JOIN "Item" ON "Item"."ID" = "GearPackItem"."itemID"
		WHERE "GearPackItem"."gearPackID" = $1
	`

	adventuringGearQuery := `
		SELECT 
			i."ID", 
			i.name, 
			i.type,
			a.description, 
			a.category, 
			a."categoryDescription", 
			i.cost, 
			i.weight,
			gpi.quantity
		FROM "AdventuringGear" a
		JOIN "Item" i ON i."ID" = a."itemID"
		JOIN "GearPackItem" gpi ON gpi."itemID" = i."ID"
		WHERE i."ID" IN (%s)
 `

	toolQuery := `
		SELECT
			i."ID",
			i.name,
			i.type,
			t.description,
			t.category,
			i.cost,
			i.weight,
			gpi.quantity
		FROM "Tool" t
		JOIN "Item" i ON i."ID" = t."itemID"
		JOIN "GearPackItem" gpi ON gpi."itemID" = i."ID"
		WHERE i."ID" IN (%s)
	 `

	for _, g := range gearPacks {
		rows, err = r.DB.Query(gearPackItemQuery, g.ID)
		if err != nil {
			log.Fatal(err)
		}

		var xItemResolver []*ItemResolver
		var adventuringIDs []string
		var toolIDs []string
		for rows.Next() {
			var id string
			var itemType string
			err = rows.Scan(
				&id,
				&itemType,
			)
			if err != nil {
				log.Fatal(err)
			}

			if itemType == "AdventuringGear" {
				adventuringIDs = append(adventuringIDs, id)
			} else if itemType == "Tool" {
				toolIDs = append(toolIDs, id)
			}
		}

		if len(adventuringIDs) > 0 {
			query := fmt.Sprintf(adventuringGearQuery, strings.Join(adventuringIDs, ", "))
			rows, err = r.DB.Query(query)
			if err != nil {
				log.Fatal(err)
			}

			for rows.Next() {
				var adventuringGear QuantifiedAdventuringGear
				var tempID int32
				rows.Scan(
					&tempID,
					&adventuringGear.Name,
					&adventuringGear.ItemType,
					&adventuringGear.Description,
					&adventuringGear.Category,
					&adventuringGear.CategoryDescription,
					&adventuringGear.Cost,
					&adventuringGear.Weight,
					&adventuringGear.Quantity,
				)

				adventuringGear.ID = Int32ToGraphqlID(tempID)
				xItemResolver = append(xItemResolver, &ItemResolver{&adventuringGear})
			}
		}

		if len(toolIDs) > 0 {
			fmt.Println("GOT A TOOL", toolIDs)
			query := fmt.Sprintf(toolQuery, strings.Join(toolIDs, ", "))
			rows, err = r.DB.Query(query)
			if err != nil {
				log.Fatal(err)
			}

			for rows.Next() {
				var tool QuantifiedTool
				var tempID int32
				rows.Scan(
					&tempID,
					&tool.Name,
					&tool.ItemType,
					&tool.Description,
					&tool.Category,
					&tool.Cost,
					&tool.Weight,
					&tool.Quantity,
				)

				tool.ID = Int32ToGraphqlID(tempID)
				xItemResolver = append(xItemResolver, &ItemResolver{&tool})
			}
		}

		g.Items = &xItemResolver
	}

	return &gearPacks
}
