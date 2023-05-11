package cli

import (
	"context"

	"github.com/ugent-library/deliver/repositories"
	"github.com/urfave/cli/v2"
)

var foldersCmd = &cli.Command{
	Name: "folders",
	Subcommands: []*cli.Command{
		expireFoldersCmd,
	},
}

var expireFoldersCmd = &cli.Command{
	Name: "expire",
	Action: func(*cli.Context) error {
		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			return err
		}
		return repo.Folders.DeleteExpired(context.Background())
	},
}
