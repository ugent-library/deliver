package repositories

import (
	"context"
	"errors"

	"github.com/ugent-library/deliver/ent"
	"github.com/ugent-library/deliver/ent/file"
	"github.com/ugent-library/deliver/models"
)

type FilesRepo struct {
	db *ent.Client
}

func (r *FilesRepo) Create(ctx context.Context, f *models.File) error {
	if err := f.Validate(); err != nil {
		return err
	}
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

func (r *FilesRepo) Exists(ctx context.Context, id string) (bool, error) {
	return r.db.File.Query().
		Where(file.IDEQ(id)).
		Exist(ctx)
}

func (r *FilesRepo) Get(ctx context.Context, id string) (*models.File, error) {
	row, err := r.db.File.Query().
		Where(file.IDEQ(id)).
		WithFolder(func(q *ent.FolderQuery) {
			q.WithSpace()
		}).
		First(ctx)
	if err != nil {
		var e *ent.NotFoundError
		if errors.As(err, &e) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return rowToFile(row), nil
}

func (r *FilesRepo) Delete(ctx context.Context, id string) error {
	err := r.db.File.
		DeleteOneID(id).
		Exec(ctx)
	return err
}

func (r *FilesRepo) AddDownload(ctx context.Context, id string) error {
	err := r.db.File.
		UpdateOneID(id).
		AddDownloads(1).
		Exec(ctx)
	return err
}

func rowToFile(row *ent.File) *models.File {
	f := &models.File{
		ID:          row.ID,
		FolderID:    row.FolderID,
		MD5:         row.Md5,
		Name:        row.Name,
		Size:        row.Size,
		ContentType: row.ContentType,
		Downloads:   row.Downloads,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
	if row.Edges.Folder != nil {
		f.Folder = rowToFolder(row.Edges.Folder)
	}
	return f
}
