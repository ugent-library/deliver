package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/deliver/objectstores"
	"github.com/ugent-library/deliver/repositories"
)

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().Bool("confirm", false, "destructive reset of all data")
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Destructive reset",
	RunE: func(cmd *cobra.Command, args []string) error {
		if confirm, _ := cmd.Flags().GetBool("confirm"); !confirm {
			return nil
		}

		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			return err
		}

		storage, err := objectstores.New(config.Storage.Backend, config.Storage.Conn)
		if err != nil {
			return err
		}

		ctx := context.TODO()

		spaces, err := repo.Spaces.GetAll(ctx)
		if err != nil {
			return err
		}

		for _, sp := range spaces {
			space, err := repo.Spaces.GetByName(ctx, sp.Name)
			if err != nil {
				return err
			}

			for _, folder := range space.Folders {
				for _, file := range folder.Files {
					if err = storage.Delete(ctx, file.ID); err != nil {
						return err
					}

					if err = repo.Files.Delete(ctx, file.ID); err != nil {
						return err
					}
				}

				if err = repo.Folders.Delete(ctx, folder.ID); err != nil {
					return err
				}
			}

			if err = repo.Spaces.Delete(ctx, space.ID); err != nil {
				return err
			}
		}

		return nil
	},
}
