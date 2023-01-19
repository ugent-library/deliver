package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/deliver/models"
)

func init() {
	rootCmd.AddCommand(foldersCmd)
	foldersCmd.AddCommand(deleteExpiredFoldersCmd)
}

var foldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "Folder commands",
}

var deleteExpiredFoldersCmd = &cobra.Command{
	Use:   "delete-expired",
	Short: "Delete all expired folders",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := models.NewRepositoryService(models.RepositoryConfig{
			DB: config.DB,
		})
		if err != nil {
			logger.Fatal(err)
		}
		if err := repo.DeleteExpiredFolders(context.TODO()); err != nil {
			logger.Fatal(err)
		}
	},
}
