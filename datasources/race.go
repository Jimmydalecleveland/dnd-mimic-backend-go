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
	q := `Select * FROM "Race" WHERE "ID" = $1`
	var race Race
	err := r.DB.QueryRow(q, args.ID).Scan(&race.ID, &race.Name)
	if err != nil {
		log.Fatal(err)
	}

	return &RaceResolver{r: &race}
}

func (r *Resolver) Races() *[]*RaceResolver {
	q := `SELECT * FROM "Race"`
	var races []*Race

	rows, err := r.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var race Race
		err = rows.Scan(&race.ID, &race.Name)
		if err != nil {
			log.Fatal(err)
		}
		races = append(races, &race)
	}

	var xRaceResolver []*RaceResolver
	for _, r := range races {
		xRaceResolver = append(xRaceResolver, &RaceResolver{r})
	}

	return &xRaceResolver
}
