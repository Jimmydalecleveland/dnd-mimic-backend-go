package datasources

import (
	"context"
	"log"

	graphql "github.com/graph-gophers/graphql-go"
)

type Background struct {
	ID                int32
	Name              string
	Description       string
	NumExtraLanguages *int32
	StartingGp        *int32
}

type BackgroundResolver struct {
	b Background
}

func (r *BackgroundResolver) ID() graphql.ID {
	return Int32ToGraphqlID(r.b.ID)
}

func (r *BackgroundResolver) Name() string {
	return r.b.Name
}

func (r *BackgroundResolver) Description() string {
	return r.b.Description
}

func (r *BackgroundResolver) NumExtraLanguages() *int32 {
	return r.b.NumExtraLanguages
}

func (r *BackgroundResolver) StartingGp() *int32 {
	return r.b.StartingGp
}

func (r *Resolver) Background(ctx context.Context, args struct{ ID int32 }) *BackgroundResolver {
	var b Background
	q := `
	SELECT * FROM "Background"
	WHERE "ID" = $1
	`

	err := r.DB.QueryRow(q, args.ID).Scan(
		&b.ID,
		&b.Name,
		&b.Description,
		&b.NumExtraLanguages,
		&b.StartingGp,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &BackgroundResolver{b}
}

func (r *Resolver) Backgrounds() *[]*BackgroundResolver {
	var backgrounds []Background

	q := `
		Select * FROM "Background"
	`

	rows, err := r.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var background Background
		err = rows.Scan(
			&background.ID,
			&background.Name,
			&background.Description,
			&background.NumExtraLanguages,
			&background.StartingGp,
		)
		if err != nil {
			log.Fatal(err)
		}
		backgrounds = append(backgrounds, background)
	}

	var xBackgroundResolver []*BackgroundResolver
	for _, b := range backgrounds {
		xBackgroundResolver = append(xBackgroundResolver, &BackgroundResolver{b})
	}
	return &xBackgroundResolver
}
