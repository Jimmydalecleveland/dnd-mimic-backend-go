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

func (Skill) TableName() string {
	return "Skill"
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

	result := r.DB.Find(&skills)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	var xSkillResolver []*SkillResolver
	for _, s := range skills {
		xSkillResolver = append(xSkillResolver, &SkillResolver{s})
	}
	return &xSkillResolver
}
