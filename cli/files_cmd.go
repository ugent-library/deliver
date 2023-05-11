package cli

import (
	"context"
	"errors"

	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstore"
	"github.com/ugent-library/deliver/repositories"
	"github.com/urfave/cli/v2"
)

var filesCmd = &cli.Command{
	Name: "files",
	Subcommands: []*cli.Command{
		gcFilesCmd,
	},
}

var gcFilesCmd = &cli.Command{
	Name: "gc",
	Action: func(*cli.Context) error {
		ctx := context.Background()

		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			return err
		}
		storage, err := objectstore.New(config.Storage.Backend, config.Storage.Conn)
		if err != nil {
			return err
		}

		var e error
		err = storage.EachID(ctx, func(id string) bool {
			_, e = repo.Files.Get(ctx, id)
			if errors.Is(e, models.ErrNotFound) {
				e = storage.Delete(ctx, id)
			}
			return e == nil
		})
		if err != nil {
			return err
		}
		return e
	},
}
