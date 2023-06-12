package repositories

import (
	"context"
	"errors"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/ugent-library/deliver/ent"
	"github.com/ugent-library/deliver/ent/user"
	"github.com/ugent-library/deliver/models"
)

type UsersRepo struct {
	db *ent.Client
}

func (r *UsersRepo) GetByRememberToken(ctx context.Context, token string) (*models.User, error) {
	row, err := r.db.User.Query().
		Where(user.RememberTokenEQ(token)).
		First(ctx)
	if err != nil {
		var e *ent.NotFoundError
		if errors.As(err, &e) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return rowToUser(row), nil

}

// TODO rewrite this when ent supports the Save method on Update; until then
// we have to do an extra select
// https://github.com/ent/ent/issues/2600
func (r *UsersRepo) CreateOrUpdate(ctx context.Context, u *models.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	token, err := models.NewRememberToken()
	if err != nil {
		return err
	}
	id, err := r.db.User.Create().
		SetUsername(u.Username).
		SetName(u.Name).
		SetEmail(u.Email).
		SetRememberToken(token).
		OnConflict(
			entsql.ConflictColumns(user.FieldUsername),
		).
		Update(func(u *ent.UserUpsert) {
			u.UpdateName()
			u.UpdateEmail()
		}).ID(ctx)
	if err != nil {
		return err
	}
	row, err := r.db.User.Get(ctx, id)
	if err != nil {
		return err
	}
	*u = *rowToUser(row)
	return nil
}

func (r *UsersRepo) RenewRememberToken(ctx context.Context, id string) error {
	newToken, err := models.NewRememberToken()
	if err != nil {
		return err
	}
	err = r.db.User.
		UpdateOneID(id).
		SetRememberToken(newToken).
		Exec(ctx)
	return err
}

func rowToUser(row *ent.User) *models.User {
	u := &models.User{
		ID:            row.ID,
		Username:      row.Username,
		Name:          row.Name,
		Email:         row.Email,
		RememberToken: row.RememberToken,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
	return u
}
