package repository

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	sqldialect "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ugent-library/deliver/ent"
	"github.com/ugent-library/deliver/ent/migrate"
)

type Repo struct {
	Users   *UsersRepo
	Spaces  *SpacesRepo
	Folders *FoldersRepo
	Files   *FilesRepo
}

func New(conn string) (*Repo, error) {
	db, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}

	driver := sqldialect.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(driver))

	err = client.Schema.Create(context.TODO(),
		migrate.WithDropIndex(true),
	)
	if err != nil {
		return nil, err
	}

	return &Repo{
		Users:   &UsersRepo{client},
		Spaces:  &SpacesRepo{client},
		Folders: &FoldersRepo{client},
		Files:   &FilesRepo{client},
	}, nil
}
