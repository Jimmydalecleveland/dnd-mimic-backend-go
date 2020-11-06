package datasources

import (
	"context"
	"fmt"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Race struct {
	ID    int32 `gorm:"column:ID"`
	Name  string
	Speed int32
}

func (Race) TableName() string {
	return "Race"
}

type RaceResolver struct {
	r *Race
}

func (r *RaceResolver) ID() graphql.ID {
	fmt.Println(r.r.ID)
	return Int32ToGraphqlID(r.r.ID)
}

func (r *RaceResolver) Name() string {
	fmt.Println(r.r.Name)
	return r.r.Name
}

func (r *Resolver) Race(ctx context.Context, args struct{ ID int32 }) *RaceResolver {
	var race Race
	result := r.DB.First(&race, args.ID)
	if result.Error != nil {
		return nil
	}

	fmt.Println(race)
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
