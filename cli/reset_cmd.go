package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/deliver/objectstores"
	"github.com/ugent-library/deliver/repositories"
)

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().Bool("force", false, "force destructive reset of all data")
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Destructive reset",
	RunE: func(cmd *cobra.Command, args []string) error {
		if force, _ := cmd.Flags().GetBool("force"); !force {
			cmd.Println("The --force flag is required to perform a destructive reset.")
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

		for _, space := range spaces {
			if err = repo.Spaces.Delete(ctx, space.ID); err != nil {
				return err
			}
		}

		err = storage.IterateID(ctx, func(id string) error {
			return storage.Delete(ctx, id)
		})

		if err != nil {
			return err
		}

		cmd.Println("Finished destructive reset.")

		return nil
	},
}
