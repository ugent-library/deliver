package repositories

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/ugent-library/deliver/ent"
	"github.com/ugent-library/deliver/ent/folder"
	"github.com/ugent-library/deliver/ent/space"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/okay"
)

type SpacesRepo struct {
	client *ent.Client
}

func (r *SpacesRepo) GetAll(ctx context.Context) ([]*models.Space, error) {
	rows, err := r.client.Space.Query().
		Order(func(s *sql.Selector) {
			s.OrderExpr(sql.ExprFunc(func(b *sql.Builder) {
				b.WriteString("UPPER(").Ident(space.FieldName).WriteString(")")
			}))
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}
	spaces := make([]*models.Space, len(rows))
	for i, row := range rows {
		spaces[i] = rowToSpace(row)
	}
	return spaces, nil
}

func (r *SpacesRepo) GetAllByUsername(ctx context.Context, username string) ([]*models.Space, error) {
	rows, err := r.client.Space.Query().
		Where(func(s *sql.Selector) {
			s.Where(sqljson.ValueContains(space.FieldAdmins, username))
		}).
		Order(ent.Asc(space.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	spaces := make([]*models.Space, len(rows))
	for i, row := range rows {
		spaces[i] = rowToSpace(row)
	}
	return spaces, nil
}

func (r *SpacesRepo) GetByName(ctx context.Context, name string) (*models.Space, error) {
	row, err := r.client.Space.Query().
		Where(space.NameEQ(name)).
		WithFolders(func(q *ent.FolderQuery) {
			q.Order(ent.Asc(folder.FieldExpiresAt))
			q.WithFiles()
		}).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return rowToSpace(row), nil
}

func (r *SpacesRepo) Create(ctx context.Context, s *models.Space) error {
	if err := s.Validate(); err != nil {
		return err
	}
	row, err := r.client.Space.Create().
		SetName(s.Name).
		SetAdmins(s.Admins).
		Save(ctx)
	if ent.IsConstraintError(err) {
		return okay.NewErrors(okay.ErrNotUnique("name"))
	}
	if err != nil {
		return err
	}
	*s = *rowToSpace(row)
	return nil
}

func (r *SpacesRepo) Update(ctx context.Context, s *models.Space) error {
	if err := s.Validate(); err != nil {
		return err
	}
	row, err := r.client.Space.UpdateOneID(s.ID).
		SetAdmins(s.Admins).
		Save(ctx)
	if err != nil {
		return err
	}
	*s = *rowToSpace(row)
	return nil
}

func (r *SpacesRepo) Delete(ctx context.Context, spaceID string) error {
	err := r.client.Space.
		DeleteOneID(spaceID).
		Exec(ctx)
	return err
}

func rowToSpace(row *ent.Space) *models.Space {
	s := &models.Space{
		ID:        row.ID,
		Name:      row.Name,
		Admins:    row.Admins,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	if row.Edges.Folders != nil {
		s.Folders = make([]*models.Folder, len(row.Edges.Folders))
		for i, r := range row.Edges.Folders {
			f := &models.Folder{
				ID:        r.ID,
				SpaceID:   r.SpaceID,
				Name:      r.Name,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
				ExpiresAt: r.ExpiresAt,
			}
			if r.Edges.Files != nil {
				f.Files = make([]*models.File, len(r.Edges.Files))
				for i, r := range r.Edges.Files {
					f.Files[i] = rowToFile(r)
				}
			}

			s.Folders[i] = f
		}
	}
	return s
}
