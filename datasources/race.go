package datasources

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Race struct {
	ID           int32 `gorm:"column:ID"`
	Name         string
	Speed        int32
	ParentRaceID int32  `gorm:"column:parentRaceID"`
	Subraces     []Race `gorm:"foreignKey:parentRaceID"`
}

// type Subrace struct {
// 	ID           int32 `gorm:"column:ID"`
// 	Name         string
// 	Speed        int32
// 	ParentRaceID int32 `gorm:"column:parentRaceID"`
// }

type RaceResolver struct {
	r *Race
}

// type SubraceResolver struct {
// 	s *Race
// }

func (Race) TableName() string {
	return "Race"
}

// func (Subrace) TableName() string {
// 	return "Race"
// }

func (r *RaceResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.r.ID)
}

func (r *RaceResolver) Name() string {
	return r.r.Name
}

func (r *RaceResolver) Subraces() *[]*RaceResolver {
	var xSubraceResolver []*RaceResolver
	for _, s := range r.r.Subraces {
		xSubraceResolver = append(xSubraceResolver, &RaceResolver{&s})
	}
	return &xSubraceResolver
}

// func (r *SubraceResolver) ID() graphql.ID {
// 	return Int32ToGraphqlID(r.s.ID)
// }

// func (r *SubraceResolver) Name() string {
// 	return r.s.Name
// }

func (r *Resolver) Race(ctx context.Context, args struct{ ID int32 }) *RaceResolver {
	var race Race
	// var raw = `
	// SELECT * FROM "Race"
	// INNER JOIN "Race" sr ON "Race"."ID" = sr."parentRaceID"
	// WHERE "Race"."ID" = ?
	// `
	// result := r.DB.Raw(raw, args.ID).Scan(&race)
	result := r.DB.First(&race, args.ID)
	if result.Error != nil {
		return nil
	}

	return &RaceResolver{r: &race}
}

func (r *Resolver) Races() *[]*RaceResolver {
	var races []*Race

	result := r.DB.Find(&races)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	var xRaceResolver []*RaceResolver
	for _, r := range races {
		xRaceResolver = append(xRaceResolver, &RaceResolver{r})
	}

	return &xRaceResolver
}

func (r *Resolver) Subrace(ctx context.Context, args struct{ ID int32 }) *RaceResolver {
	var subrace Race
	result := r.DB.Preload("Race").First(&subrace, args.ID)
	if result.Error != nil {
		return nil
	}

	return &RaceResolver{r: &subrace}
}
