package datasources

import (
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Skill struct {
	ID      int32
	Name    string
	Ability string
}

type SkillResolver struct {
	s *Skill
}

func (r *SkillResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.s.ID)
}

func (r *SkillResolver) Name() string {
	return r.s.Name
}

func (r *SkillResolver) Ability() string {
	return r.s.Ability
}

func (r *Resolver) Skills() *[]*SkillResolver {
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
		err = rows.Scan(&singleSkill.ID, &singleSkill.Name, &singleSkill.Ability)
		if err != nil {
			log.Fatal(err)
		}
		skills = append(skills, &singleSkill)
	}

	var xSkillResolver []*SkillResolver
	for _, s := range skills {
		xSkillResolver = append(xSkillResolver, &SkillResolver{s})
	}
	return &xSkillResolver
}
