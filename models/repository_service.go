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
	db, err := sql.Open("pgx", c.DB)
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
	spaces := make([]*Space, len(rows))
	for i, row := range rows {
		spaces[i] = rowToSpace(row)
	}
	return spaces, nil
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
	return rowToSpace(row), nil
}

func (r *repositoryService) CreateSpace(ctx context.Context, s *Space) error {
	row, err := r.db.Space.Create().SetName(s.Name).Save(ctx)
	if err != nil {
		return err
	}
	*s = *rowToSpace(row)
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
	return rowToFolder(row), nil
}

func (r *repositoryService) CreateFolder(ctx context.Context, f *Folder) error {
	row, err := r.db.Folder.Create().
		SetSpaceID(f.SpaceID).
		SetName(f.Name).
		Save(ctx)
	if err != nil {
		return err
	}
	*f = *rowToFolder(row)
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
		SetMd5(f.MD5).
		SetName(f.Name).
		SetContentType(f.ContentType).
		SetSize(f.Size).
		Save(ctx)
	if err != nil {
		return err
	}
	*f = *rowToFile(row)
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
	return rowToFile(row), nil
}

func (r *repositoryService) DeleteFile(ctx context.Context, fileID string) error {
	err := r.db.File.
		DeleteOneID(fileID).
		Exec(ctx)
	return err
}

func rowToSpace(row *ent.Space) *Space {
	s := &Space{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	if row.Edges.Folders != nil {
		s.Folders = make([]*Folder, len(row.Edges.Folders))
		for i, r := range row.Edges.Folders {
			s.Folders[i] = rowToFolder(r)
		}
	}
	return s
}

func rowToFolder(row *ent.Folder) *Folder {
	f := &Folder{
		ID:        row.ID,
		SpaceID:   row.SpaceID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		ExpiresAt: row.ExpiresAt,
	}
	if row.Edges.Space != nil {
		f.Space = rowToSpace(row.Edges.Space)
	}
	if row.Edges.Files != nil {
		f.Files = make([]*File, len(row.Edges.Files))
		for i, r := range row.Edges.Files {
			f.Files[i] = rowToFile(r)
		}
	}
	return f
}

func rowToFile(row *ent.File) *File {
	f := &File{
		ID:          row.ID,
		FolderID:    row.FolderID,
		MD5:         row.Md5,
		Name:        row.Name,
		Size:        row.Size,
		ContentType: row.ContentType,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
	if row.Edges.Folder != nil {
		f.Folder = rowToFolder(row.Edges.Folder)
	}
	return f
}
