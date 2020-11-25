package datasources

import (
	"log"

	"github.com/graph-gophers/graphql-go"
)

type Class struct {
	ID                       int32
	Name                     string
	HitDice                  string
	NumSkillProficiencies    int32
	SavingThrowProficiencies []uint8
}

type ClassResolver struct {
	c Class
}

func (r *ClassResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.c.ID)
}

func (r ClassResolver) Name() string {
	return r.c.Name
}

func (r ClassResolver) HitDice() string {
	return r.c.HitDice
}

func (r ClassResolver) NumSkillProficiencies() int32 {
	return r.c.NumSkillProficiencies
}

func (r *Resolver) Classes() *[]*ClassResolver {
	var classes []Class
	q := `
		SELECT * FROM "CharClass"	
	`

	rows, err := r.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var class Class
		err = rows.Scan(
			&class.ID,
			&class.Name,
			&class.HitDice,
			&class.NumSkillProficiencies,
			&class.SavingThrowProficiencies,
		)
		if err != nil {
			log.Fatal(err)
		}
		classes = append(classes, class)
	}

	var xClassResolver []*ClassResolver
	for _, c := range classes {
		xClassResolver = append(xClassResolver, &ClassResolver{c})
	}

	return &xClassResolver
}
