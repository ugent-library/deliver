package cmd

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
	Use:   "folders",
	Short: "Folder commands",
}

var expireFoldersCmd = &cobra.Command{
	Use:   "expire",
	Short: "Delete all expired folders",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			logger.Fatal(err)
		}
		if err := repo.Folders.DeleteExpired(context.Background()); err != nil {
			logger.Fatal(err)
		}
	},
}
