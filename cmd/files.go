package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	"github.com/ugent-library/deliver/models"
)

func init() {
	rootCmd.AddCommand(filesCmd)
	foldersCmd.AddCommand(gcFilesCmd)
}

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "File commands",
}

var gcFilesCmd = &cobra.Command{
	Use:   "gc",
	Short: "Garbage collect orphaned files",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()

		repoService, err := models.NewRepositoryService(models.RepositoryConfig{
			DB: config.DB,
		})
		if err != nil {
			logger.Fatal(err)
		}
		fileService, err := models.NewFileService(models.FileConfig{
			S3ID:     config.S3.ID,
			S3Secret: config.S3.Secret,
			S3Bucket: config.S3.Bucket,
			S3Region: config.S3.Region,
		})
		if err != nil {
			logger.Fatal(err)
		}

		err = fileService.EachID(ctx, func(id string) bool {
			_, err = repoService.FileByID(ctx, id)
			if errors.Is(err, models.ErrNotFound) {
				err = fileService.Delete(ctx, id)
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
