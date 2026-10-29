package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/deliver/repositories"
)

func init() {
	rootCmd.AddCommand(foldersCmd)
	foldersCmd.AddCommand(expireFoldersCmd)
}

var foldersCmd = &cobra.Command{
	Use: "folders",
}

var expireFoldersCmd = &cobra.Command{
	Use: "expire",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			return err
		}

		logger.Info("deleting expired folders")

		n, err := repo.Folders.DeleteExpired(context.Background())
		if err != nil {
			logger.Errorf("could not delete expired folders: %s", err)
			return err
		}

		logger.Infof("deleted %d expired folders", n)

		return nil
	},
}
