package models

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	entdialect "entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/ugent-library/dilliver/ent"
	entmigrate "github.com/ugent-library/dilliver/ent/migrate"
)

type Space struct {
	ID   string
	Name string
}

type RepositoryService interface {
	Spaces(context.Context) ([]*Space, error)
	CreateSpace(context.Context, *Space) error
}

func NewRepositoryService(c Config) (RepositoryService, error) {
	db, err := sql.Open("pgx", c.DatabaseURL)
	if err != nil {
		return nil, err
	}

	driver := entsql.OpenDB(entdialect.Postgres, db)
	client := ent.NewClient(ent.Driver(driver))

	err = client.Schema.Create(context.Background(),
		entmigrate.WithDropIndex(true),
	)
	if err != nil {
		return nil, err
	}

	return &repository{
		db: client,
	}, nil
}

type repository struct {
	db *ent.Client
}

func (r *repository) Spaces(ctx context.Context) ([]*Space, error) {
	rows, err := r.db.Space.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	objs := make([]*Space, len(rows))
	for i, row := range rows {
		objs[i] = &Space{
			ID:   row.ID,
			Name: row.Name,
		}
	}
	return objs, nil
}

func (r *repository) CreateSpace(ctx context.Context, s *Space) error {
	row, err := r.db.Space.Create().SetName(s.Name).Save(ctx)
	if err != nil {
		return err
	}
	s.ID = row.ID
	return nil
}
