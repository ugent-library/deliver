package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/deliver/models"
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
		repoService, err := models.NewRepositoryService(models.RepositoryConfig{
			DB: config.DB,
		})
		if err != nil {
			logger.Fatal(err)
		}
		if err := repoService.DeleteExpiredFolders(context.TODO()); err != nil {
			logger.Fatal(err)
		}
	},
}
