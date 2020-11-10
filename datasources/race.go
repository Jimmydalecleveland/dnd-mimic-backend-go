package datasources

import (
	"context"
	"database/sql"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Race struct {
	ID           int32
	Name         string
	Speed        int32
	ParentRaceID int32
}

type RaceResolver struct {
	r Race
	s []Race
}

type SubraceResolver struct {
	s Race
}

func (Race) TableName() string {
	return "Race"
}

func (r *RaceResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.r.ID)
}

func (r *RaceResolver) Name() string {
	return r.r.Name
}

func (r *RaceResolver) Subraces() *[]*SubraceResolver {
	var xSubraceResolver []*SubraceResolver
	for _, s := range r.s {
		xSubraceResolver = append(xSubraceResolver, &SubraceResolver{s})
	}

	return &xSubraceResolver
}

func (r *SubraceResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.s.ID)
}

func (r *SubraceResolver) Name() string {
	return r.s.Name
}

func queryRace(db *sql.DB, id int32) Race {
	q := `Select "Race"."ID", "Race".name FROM "Race" WHERE "ID" = $1`
	var race Race
	err := db.QueryRow(q, id).Scan(&race.ID, &race.Name)
	if err != nil {
		log.Fatal(err)
	}

	return race
}

func querySubraces(db *sql.DB, id int32) []Race {
	q2 := `
		SELECT subrace.* FROM "Race"
		JOIN "Race" subrace ON "Race"."ID" = subrace."parentRaceID"
		WHERE "Race"."ID" = $1
	`
	var subraces []Race
	rows, err := db.Query(q2, id)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var singleRace Race
		err = rows.Scan(&singleRace.ID, &singleRace.Name, &singleRace.ParentRaceID, &singleRace.Speed)
		if err != nil {
			log.Fatal(err)
		}
		subraces = append(subraces, singleRace)
	}

	return subraces
}

func (r *Resolver) Race(ctx context.Context, args struct{ ID int32 }) *RaceResolver {
	race := queryRace(r.DB, args.ID)
	subraces := querySubraces(r.DB, args.ID)

	return &RaceResolver{r: race, s: subraces}
}

// func (r *Resolver) Races() *[]*RaceResolver {
// 	q := `SELECT * FROM "Race"`
// 	var races []*Race

// 	rows, err := r.DB.Query(q)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for rows.Next() {
// 		var race Race
// 		err = rows.Scan(&race.ID, &race.Name)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		races = append(races, &race)
// 	}

// 	var xRaceResolver []*RaceResolver
// 	for _, r := range races {
// 		xRaceResolver = append(xRaceResolver, &RaceResolver{r})
// 	}

// 	return &xRaceResolver
// }
