package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/oklog/ulid/v2"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstores"
	"github.com/ugent-library/deliver/repositories"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.Flags().Bool("force", false, "force seeding the database")
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the application with dummy data",
	RunE: func(cmd *cobra.Command, args []string) error {
		// setup services
		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			return err
		}

		storage, err := objectstores.New(config.Storage.Backend, config.Storage.Conn)
		if err != nil {
			return err
		}

		ctx := context.TODO()

		if force, _ := cmd.Flags().GetBool("force"); !force {
			spaces, err := repo.Spaces.GetAll(ctx)
			if err != nil {
				return err
			}

			if len(spaces) > 0 {
				fmt.Println("The database is not empty. Not seeding.")
				return nil
			}
		}

		// Create a 'deliver' user
		var user models.User
		gofakeit.Struct(&user)

		if err := repo.Users.CreateOrUpdate(ctx, &user); err != nil {
			return err
		}

		// Create a space
		for i := 0; i < 5; i++ {
			var space models.Space
			gofakeit.Struct(&space)

			if err := repo.Spaces.Create(ctx, &space); err != nil {
				return err
			}

			for j := 0; j < 5; j++ {
				// Create a folder
				var folder models.Folder
				gofakeit.Struct(&folder)
				folder.SpaceID = space.ID

				if err := repo.Folders.Create(ctx, &folder); err != nil {
					return err
				}

				for k := 0; k < 5; k++ {
					// Create a file
					var file models.File
					gofakeit.Struct(&file)

					file.FolderID = folder.ID
					file.ID = ulid.Make().String()

					ba := gofakeit.ImageJpeg(1024, 1024)
					br := bytes.NewReader(ba)
					rc := io.NopCloser(br)

					md5, err := storage.Add(ctx, file.ID, rc)
					if err != nil {
						return err
					}

					file.Size = int64(len(ba))
					file.MD5 = md5
					file.ContentType = "image/jpeg"

					if err := repo.Files.Create(ctx, &file); err != nil {
						return err
					}
				}
			}
		}

		fmt.Println("Done.")

		return nil
	},
}
