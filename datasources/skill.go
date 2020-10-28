package datasources

import (
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type skill struct {
	ID      int32
	name    string
	ability string
}

type SkillResolver struct {
	s *skill
}

func (r *SkillResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.s.ID)
}

func (r *SkillResolver) Name() string {
	return r.s.name
}

func (r *SkillResolver) Ability() string {
	return r.s.ability
}

func (r *Resolver) Skills() *[]*SkillResolver {
	var skills []*skill
	var err error
	skillQuery := `
		Select "ID", name, ability FROM "Skill"
	`

	rows, err := r.DB.Query(skillQuery)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var singleSkill skill
		err = rows.Scan(&singleSkill.ID, &singleSkill.name, &singleSkill.ability)
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
