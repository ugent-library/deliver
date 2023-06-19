package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/deliver/objectstore"
	repository "github.com/ugent-library/deliver/repository"
)

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.AddCommand(gcFilesCmd)
}

var filesCmd = &cobra.Command{
	Use: "files",
}

var gcFilesCmd = &cobra.Command{
	Use: "gc",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		repo, err := repository.New(config.Repo.Conn)
		if err != nil {
			return err
		}
		storage, err := objectstore.New(config.Storage.Backend, config.Storage.Conn)
		if err != nil {
			return err
		}

		return storage.IterateID(ctx, func(id string) error {
			exists, err := repo.Files.Exists(ctx, id)
			if err != nil {
				return err
			}
			if !exists {
				if err = storage.Delete(ctx, id); err != nil {
					return err
				}
			}
			return nil
		})
	},
}
