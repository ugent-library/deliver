package cmd

import (
	"encoding/gob"
	"html/template"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	c "github.com/ugent-library/dilliver/controllers"
	"github.com/ugent-library/dilliver/mix"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/view"
)

func init() {
	rootCmd.AddCommand(appCmd)
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Start the web app server",
	Run: func(cmd *cobra.Command, args []string) {
		// setup services
		services, err := models.NewServices(models.Config{
			DatabaseURL:       viper.GetString("database_url"),
			S3URL:             viper.GetString("s3_url"),
			S3AccessKeyID:     viper.GetString("s3_access_key_id"),
			S3SecretAccessKey: viper.GetString("s3_secret_access_key"),
			S3Bucket:          viper.GetString("s3_bucket"),
		})
		if err != nil {
			logger.Fatal(err)
		}

		// setup assets
		assets, err := mix.New(mix.Config{
			ManifestFile: "static/mix-manifest.json",
			PublicPath:   "/static/",
		})
		if err != nil {
			logger.Fatal(err)
		}

		// setup router
		r := mux.NewRouter()
		r.StrictSlash(true)
		r.UseEncodedPath()
		r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

		// setup views
		view.FuncMap = template.FuncMap{
			"assetPath": assets.AssetPath,
		}

		// setup sessions
		sessionName := viper.GetString("session_name")
		sessionStore := sessions.NewCookieStore([]byte(viper.GetString("session_secret")))
		sessionStore.MaxAge(viper.GetInt("session_max_age"))
		sessionStore.Options.Path = "/"
		sessionStore.Options.HttpOnly = true
		sessionStore.Options.Secure = viper.GetBool("production")
		// register Flash as a gob Type so CookieStore can serialize it
		gob.Register(c.Flash{})

		// controllers
		pages := c.NewPages()
		spaces := c.NewSpaces(services.Repository)
		folders := c.NewFolders(services.Repository, services.File)

		// request context wrapper
		wrap := c.Wrapper(c.Config{
			Log:          logger,
			SessionStore: sessionStore,
			SessionName:  sessionName,
			Router:       r,
		})

		// routes
		r.NotFoundHandler = wrap(pages.NotFound)
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		r.HandleFunc("/", wrap(pages.Home)).Methods("GET").Name("home")
		r.HandleFunc("/spaces", wrap(spaces.List)).Methods("GET").Name("spaces")
		r.HandleFunc("/spaces", wrap(spaces.Create)).Methods("POST").Name("create_space")
		r.HandleFunc("/spaces/{spaceID}", wrap(spaces.Show)).Methods("GET").Name("space")
		r.HandleFunc("/spaces/{spaceID}/folders", wrap(folders.Create)).Methods("POST").Name("create_folder")
		r.HandleFunc("/folders/{folderID}", wrap(folders.Show)).Methods("GET").Name("folder")
		r.HandleFunc("/folders/{folderID}/files", wrap(folders.UploadFile)).Methods("POST").Name("upload_file")
		r.HandleFunc("/folders/{folderID}/files/{fileID}", wrap(folders.DownloadFile)).Methods("GET").Name("download_file")

		// start server
		if err = http.ListenAndServe(viper.GetString("app_addr"), r); err != nil {
			logger.Fatal(err)
		}
	},
}
