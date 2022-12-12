package models

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	entdialect "entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/ugent-library/dilliver/ent"
	"github.com/ugent-library/dilliver/ent/folder"
	entmigrate "github.com/ugent-library/dilliver/ent/migrate"
)

type Space = ent.Space
type Folder = ent.Folder
type File = ent.File

type RepositoryService interface {
	Spaces(context.Context) ([]*Space, error)
	CreateSpace(context.Context, *Space) error
	Folders(context.Context, string) ([]*Folder, error)
	Folder(context.Context, string) (*Folder, error)
	CreateFolder(context.Context, *Folder) error
	CreateFile(context.Context, *File) error
}

func NewRepositoryService(c Config) (RepositoryService, error) {
	db, err := sql.Open("pgx", c.DatabaseURL)
	if err != nil {
		return nil, err
	}

	driver := entsql.OpenDB(entdialect.Postgres, db)
	client := ent.NewClient(ent.Driver(driver))

	err = client.Schema.Create(context.TODO(),
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

func (r *repositoryService) Folder(ctx context.Context, folderID string) (*Folder, error) {
	folder, err := r.db.Folder.Query().
		Where(folder.IDEQ(folderID)).
		WithFiles().
		First(ctx)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (r *repositoryService) CreateFolder(ctx context.Context, f *Folder) error {
	folder, err := r.db.Folder.Create().
		SetSpaceID(f.SpaceID).
		SetName(f.Name).
		Save(ctx)
	if err != nil {
		return err
	}
	*f = *folder
	return nil
}

func (r *repositoryService) CreateFile(ctx context.Context, f *File) error {
	file, err := r.db.File.Create().
		SetID(f.ID).
		SetFolderID(f.FolderID).
		SetMd5(f.Md5).
		SetName(f.Name).
		SetContentType(f.ContentType).
		SetSize(f.Size).
		Save(ctx)
	if err != nil {
		return err
	}
	*f = *file
	return nil
}
