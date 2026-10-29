package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/deliver/objectstores"
	"github.com/ugent-library/deliver/repositories"
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

		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			return err
		}
		storage, err := objectstores.New(config.Storage.Backend, config.Storage.Conn)
		if err != nil {
			return err
		}

		logger.Info("starting file garbage collection")

		var checked, deleted int

		err = storage.IterateID(ctx, func(id string) error {
			checked++
			exists, err := repo.Files.Exists(ctx, id)
			if err != nil {
				logger.Errorf("could not check if file %s exists: %s", id, err)
				return err
			}
			if !exists {
				logger.Infof("deleting orphaned file %s", id)
				if err = storage.Delete(ctx, id); err != nil {
					logger.Errorf("could not delete file %s: %s", id, err)
					return err
				}
				deleted++
			}
			return nil
		})
		if err != nil {
			logger.Errorf("file garbage collection aborted after checking %d files (%d deleted): %s", checked, deleted, err)
			return err
		}

		logger.Infof("finished file garbage collection: checked %d files, deleted %d", checked, deleted)

		return nil
	},
}
