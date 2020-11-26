package datasources

import (
	"context"
	"database/sql"
	"log"

	"github.com/graph-gophers/graphql-go"
	"github.com/lib/pq"
)

type Class struct {
	ID                       int32
	Name                     string
	HitDice                  string
	NumSkillProficiencies    int32
	SavingThrowProficiencies []string
}

type ClassResolver struct {
	c Class
}

func (r ClassResolver) ID() graphql.ID {
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

func (r ClassResolver) SavingThrowProficiencies() *[]string {
	return &r.c.SavingThrowProficiencies
}

func queryClass(db *sql.DB, id int32) Class {
	var c Class
	q := `
		SELECT * FROM "CharClass"
		WHERE "ID" = $1
	`
	err := db.QueryRow(q, id).Scan(
		&c.ID,
		&c.Name,
		&c.HitDice,
		&c.NumSkillProficiencies,
		pq.Array(&c.SavingThrowProficiencies),
	)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func (r *Resolver) Class(ctx context.Context, args struct{ ID int32 }) *ClassResolver {
	c := queryClass(r.DB, args.ID)
	return &ClassResolver{c}
}

func (r *Resolver) Classes() *[]ClassResolver {
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
			pq.Array(&class.SavingThrowProficiencies),
		)
		if err != nil {
			log.Fatal(err)
		}
		classes = append(classes, class)
	}

	var xClassResolver []ClassResolver
	for _, c := range classes {
		xClassResolver = append(xClassResolver, ClassResolver{c})
	}

	return &xClassResolver
}
