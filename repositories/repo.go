package repositories

import (
	"database/sql"

	"entgo.io/ent/dialect"
	sqldialect "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ugent-library/deliver/ent"
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

	return &Repo{
		Users:   &UsersRepo{client},
		Spaces:  &SpacesRepo{client},
		Folders: &FoldersRepo{client},
		Files:   &FilesRepo{client},
	}, nil
}
