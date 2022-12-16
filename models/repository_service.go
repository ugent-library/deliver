package models

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"

	entdialect "entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/ugent-library/dilliver/ent"
	"github.com/ugent-library/dilliver/ent/file"
	"github.com/ugent-library/dilliver/ent/folder"
	entmigrate "github.com/ugent-library/dilliver/ent/migrate"
	"github.com/ugent-library/dilliver/ent/space"
)

var ErrNotFound = errors.New("not found")

type Space = ent.Space
type Folder = ent.Folder
type File = ent.File

type RepositoryService interface {
	Spaces(context.Context) ([]*Space, error)
	Space(context.Context, string) (*Space, error)
	CreateSpace(context.Context, *Space) error
	Folder(context.Context, string) (*Folder, error)
	CreateFolder(context.Context, *Folder) error
	DeleteFolder(context.Context, string) error
	File(context.Context, string) (*File, error)
	CreateFile(context.Context, *File) error
	DeleteFile(context.Context, string) error
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
	rows, err := r.db.Space.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *repositoryService) Space(ctx context.Context, spaceID string) (*Space, error) {
	row, err := r.db.Space.Query().
		Where(space.IDEQ(spaceID)).
		WithFolders().
		First(ctx)
	if err != nil {
		var e *ent.NotFoundError
		if errors.As(err, &e) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (r *repositoryService) CreateSpace(ctx context.Context, s *Space) error {
	row, err := r.db.Space.Create().SetName(s.Name).Save(ctx)
	if err != nil {
		return err
	}
	*s = *row
	return nil
}

func (r *repositoryService) Folder(ctx context.Context, folderID string) (*Folder, error) {
	row, err := r.db.Folder.Query().
		Where(folder.IDEQ(folderID)).
		WithSpace().
		WithFiles().
		First(ctx)
	if err != nil {
		var e *ent.NotFoundError
		if errors.As(err, &e) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (r *repositoryService) CreateFolder(ctx context.Context, f *Folder) error {
	row, err := r.db.Folder.Create().
		SetSpaceID(f.SpaceID).
		SetName(f.Name).
		Save(ctx)
	if err != nil {
		return err
	}
	*f = *row
	return nil
}

func (r *repositoryService) DeleteFolder(ctx context.Context, folderID string) error {
	err := r.db.Folder.
		DeleteOneID(folderID).
		Exec(ctx)
	return err
}

func (r *repositoryService) CreateFile(ctx context.Context, f *File) error {
	row, err := r.db.File.Create().
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
	*f = *row
	return nil
}

func (r *repositoryService) File(ctx context.Context, fileID string) (*File, error) {
	row, err := r.db.File.Query().
		Where(file.IDEQ(fileID)).
		WithFolder().
		First(ctx)
	if err != nil {
		var e *ent.NotFoundError
		if errors.As(err, &e) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return row, nil
}

func (r *repositoryService) DeleteFile(ctx context.Context, fileID string) error {
	err := r.db.File.
		DeleteOneID(fileID).
		Exec(ctx)
	return err
}
