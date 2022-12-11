package models

import (
	"context"
	"database/sql"
	"io"

	_ "github.com/jackc/pgx/v5/stdlib"

	entdialect "entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/ugent-library/dilliver/ent"
	entmigrate "github.com/ugent-library/dilliver/ent/migrate"
)

type Space = ent.Space
type Folder = ent.Folder
type File = ent.File

type RepositoryService interface {
	Spaces(context.Context) ([]*Space, error)
	CreateSpace(context.Context, *Space) error
	Folders(context.Context, string) ([]*Folder, error)
	CreateFolder(context.Context, *Folder) error
	CreateFile(context.Context, *File, io.Reader) error
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

	return &repositoryService{
		db: client,
	}, nil
}

type repositoryService struct {
	db *ent.Client
}

func (r *repositoryService) Spaces(ctx context.Context) ([]*Space, error) {
	spaces, err := r.db.Space.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return spaces, nil
}

func (r *repositoryService) CreateSpace(ctx context.Context, s *Space) error {
	space, err := r.db.Space.Create().SetName(s.Name).Save(ctx)
	if err != nil {
		return err
	}
	*s = *space
	return nil
}

func (r *repositoryService) Folders(ctx context.Context, spaceID string) ([]*Folder, error) {
	folders, err := r.db.Folder.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func (r *repositoryService) CreateFolder(ctx context.Context, f *Folder) error {
	folder, err := r.db.Folder.Create().
		SetName(f.Name).
		SetSpaceID(f.SpaceID).
		Save(ctx)
	if err != nil {
		return err
	}
	*f = *folder
	return nil
}

func (r *repositoryService) CreateFile(ctx context.Context, f *File, b io.Reader) error {
	return nil
}
