package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstore"
	"github.com/ugent-library/deliver/repositories"
)

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.AddCommand(gcFilesCmd)
}

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "File commands",
}

var gcFilesCmd = &cobra.Command{
	Use:   "gc",
	Short: "Garbage collect orphaned files",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			logger.Fatal(err)
		}
		storage, err := objectstore.New(config.Storage.Backend, config.Storage.Conn)
		if err != nil {
			logger.Fatal(err)
		}

		err = storage.EachID(ctx, func(id string) bool {
			_, err = repo.Files.Get(ctx, id)
			if errors.Is(err, models.ErrNotFound) {
				err = storage.Delete(ctx, id)
			}
			if err != nil {
				logger.Fatal(err)
			}
			return true
		})
		if err != nil {
			logger.Fatal(err)
		}
	},
}
