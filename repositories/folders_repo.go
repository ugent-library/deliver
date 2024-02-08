package repositories

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/ugent-library/deliver/ent"
	"github.com/ugent-library/deliver/ent/file"
	"github.com/ugent-library/deliver/ent/folder"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/okay"
)

type FoldersRepo struct {
	client *ent.Client
}

func (r *FoldersRepo) Get(ctx context.Context, id string) (*models.Folder, error) {
	row, err := r.client.Folder.Query().
		Where(folder.IDEQ(id)).
		WithSpace().
		WithFiles(func(q *ent.FileQuery) {
			q.Order(ent.Asc(file.FieldName))
		}).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return rowToFolder(row), nil
}

func (r *FoldersRepo) GetBySpace(ctx context.Context, space *models.Space, q string) ([]*models.Folder, error) {
	query := r.client.Folder.Query().
		Where(folder.SpaceIDEQ(space.ID))

	if q != "" {
		query = query.Where(func(s *sql.Selector) {
			s.Where(sql.Like(folder.FieldName, "%"+q+"%"))
		})
	}
	rows, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	folders := make([]*models.Folder, len(rows))
	for i, row := range rows {
		folders[i] = rowToFolder(row)
	}
	return folders, nil
}

func (r *FoldersRepo) Create(ctx context.Context, f *models.Folder) error {
	if err := f.Validate(); err != nil {
		return err
	}
	row, err := r.client.Folder.Create().
		SetSpaceID(f.SpaceID).
		SetName(f.Name).
		SetExpiresAt(f.ExpiresAt).
		Save(ctx)
	if ent.IsConstraintError(err) {
		return okay.NewErrors(okay.ErrNotUnique("name"))
	}
	if err != nil {
		return err
	}
	*f = *rowToFolder(row)
	return nil
}

func (r *FoldersRepo) Update(ctx context.Context, f *models.Folder) error {
	if err := f.Validate(); err != nil {
		return err
	}
	row, err := r.client.Folder.UpdateOneID(f.ID).
		SetName(f.Name).
		SetExpiresAt(f.ExpiresAt).
		Save(ctx)
	if ent.IsConstraintError(err) {
		return okay.NewErrors(okay.ErrNotUnique("name"))
	}
	if err != nil {
		return err
	}
	*f = *rowToFolder(row)
	return nil
}

func (r *FoldersRepo) Delete(ctx context.Context, folderID string) error {
	err := r.client.Folder.
		DeleteOneID(folderID).
		Exec(ctx)
	return err
}

func (r *FoldersRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.client.Folder.
		Delete().
		Where(folder.ExpiresAtLT(time.Now())).
		Exec(ctx)
	return err
}

func rowToFolder(row *ent.Folder) *models.Folder {
	f := &models.Folder{
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
		f.Files = make([]*models.File, len(row.Edges.Files))
		for i, r := range row.Edges.Files {
			f.Files[i] = rowToFile(r)
		}
	}
	return f
}
