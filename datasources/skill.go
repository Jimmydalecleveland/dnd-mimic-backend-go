package datasources

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Skill struct {
	ID      graphql.ID
	Name    string
	Ability string
}

func (r *Resolver) Skill(ctx context.Context, args struct{ ID int32 }) *Skill {
	q := `
		SELECT "ID", name, ability FROM "Skill"
		WHERE "ID" = $1
	`

	var skill Skill
	var tempID int32
	err := r.DB.QueryRow(q, args.ID).Scan(
		&tempID,
		&skill.Name,
		&skill.Ability,
	)
	if err != nil {
		log.Println("Error received during 'Skill' query -> ", err)
		return nil
	}
	skill.ID = Int32ToGraphqlID(tempID)

	return &skill
}

func (r *Resolver) Skills() *[]*Skill {
	var skills []*Skill

	var err error
	skillQuery := `
		Select "ID", name, ability FROM "Skill"
	`

	rows, err := r.DB.Query(skillQuery)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var singleSkill Skill
		var tempID int32
		err = rows.Scan(
			&tempID,
			&singleSkill.Name,
			&singleSkill.Ability,
		)
		if err != nil {
			log.Fatal(err)
		}
		singleSkill.ID = Int32ToGraphqlID(tempID)
		skills = append(skills, &singleSkill)
	}

	return &skills
}
