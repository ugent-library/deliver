package cmd

import (
	"context"
	"encoding/gob"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/spf13/cobra"
	"github.com/ugent-library/dilliver/autosession"
	c "github.com/ugent-library/dilliver/controllers"
	"github.com/ugent-library/dilliver/handler"
	"github.com/ugent-library/dilliver/middleware"
	"github.com/ugent-library/dilliver/mix"
	"github.com/ugent-library/dilliver/models"
	"github.com/ugent-library/dilliver/oidc"
	"github.com/ugent-library/dilliver/ulid"
	"github.com/ugent-library/dilliver/view"
	"github.com/ugent-library/dilliver/zaphttp"
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
			DB:       config.DB,
			S3URL:    config.S3.URL,
			S3ID:     config.S3.ID,
			S3Secret: config.S3.Secret,
			S3Bucket: config.S3.Bucket,
			S3Region: config.S3.Region,
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

		// setup views
		view.DefaultConfig.Funcs = template.FuncMap{
			"assetPath": assets.AssetPath,
		}

		// setup sessions
		sessionName := config.Session.Name
		sessionStore := sessions.NewCookieStore([]byte(config.Session.Secret))
		sessionStore.MaxAge(config.Session.MaxAge)
		sessionStore.Options.Path = "/"
		sessionStore.Options.HttpOnly = true
		sessionStore.Options.Secure = config.Production
		// register types so CookieStore can serialize it
		gob.Register(&models.User{})
		gob.Register(c.Flash{})

		// setup auth
		oidcAuth, err := oidc.NewAuth(context.TODO(), oidc.Config{
			URL:          config.Oidc.URL,
			ClientID:     config.Oidc.ID,
			ClientSecret: config.Oidc.Secret,
			RedirectURL:  config.Oidc.RedirectURL,
			CookieName:   config.Session.Name + ".state",
			CookieSecret: []byte(config.Session.Secret),
		})
		if err != nil {
			logger.Fatal(err)
		}

		// setup router
		r := mux.NewRouter()
		r.StrictSlash(true)
		r.UseEncodedPath()
		r.Use(handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			// TODO
			// handlers.RecoveryLogger(&recoveryLogger{logger}),
		))
		r.Use(csrf.Protect(
			[]byte(config.Session.Secret),
			csrf.CookieName(config.Session.Name+".csrf"),
			csrf.Path("/"),
			csrf.Secure(config.Production),
			csrf.SameSite(csrf.SameSiteStrictMode),
			csrf.FieldName("csrf_token"),
		))
		r.Use(autosession.Enable(
			autosession.GorillaSession(sessionStore, sessionName),
		))

		// controllers
		errs := c.NewErrors()
		auth := c.NewAuth(oidcAuth)
		pages := c.NewPages()
		spaces := c.NewSpaces(services.Repository)
		folders := c.NewFolders(services.Repository, services.File)
		files := c.NewFiles(services.Repository, services.File)

		// request context wrapper
		wrap := handler.Config[c.Var]{
			Log:    logger,
			Router: r,
			Before: []func(c.Ctx) error{
				c.LoadSession,
			},
			ErrorHandler: errs.HandleError,
		}.Wrap

		// routes
		r.NotFoundHandler = wrap(errs.NotFound)
		// TODO don't apply all middleware to static file server
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		r.Handle("/", wrap(pages.Home)).Methods("GET").Name("home")
		r.Handle("/auth/callback", wrap(auth.Callback)).Methods("GET")
		r.Handle("/logout", wrap(auth.Logout)).Methods("GET").Name("logout")
		r.Handle("/login", wrap(auth.Login)).Methods("GET").Name("login")
		r.Handle("/spaces", wrap(c.RequireUser, spaces.List)).Methods("GET").Name("spaces")
		r.Handle("/spaces", wrap(c.RequireUser, spaces.Create)).Methods("POST").Name("create_space")
		r.Handle("/spaces/{spaceID}", wrap(c.RequireUser, spaces.Show)).Methods("GET").Name("space")
		r.Handle("/spaces/{spaceID}/folders", wrap(c.RequireUser, folders.Create)).Methods("POST").Name("create_folder")
		r.Handle("/folders/{folderID}", wrap(c.RequireUser, folders.Show)).Methods("GET").Name("folder")
		r.Handle("/folders/{folderID}", wrap(c.RequireUser, folders.Delete)).Methods("DELETE").Name("delete_folder")
		r.Handle("/folders/{folderID}/files", wrap(c.RequireUser, folders.UploadFile)).Methods("POST").Name("upload_file")
		r.Handle("/files/{fileID}", wrap(files.Download)).Methods("GET").Name("download_file")
		r.Handle("/files/{fileID}", wrap(c.RequireUser, files.Delete)).Methods("DELETE").Name("delete_file")

		// apply method overwrite and logging handlers before request reaches the router
		// TODO Chain function to make this more readable
		var handler http.Handler = r
		handler = zaphttp.LogRequests(handler)
		handler = zaphttp.SetLogger(logger.Desugar())(handler)
		handler = middleware.SetRequestID(ulid.MustGenerate)(handler)
		if config.Production {
			handler = handlers.ProxyHeaders(handler)
		}
		handler = handlers.HTTPMethodOverrideHandler(handler)

		// start server
		// TODO timeouts, graceful shutdown
		if err = http.ListenAndServe(config.Addr, handler); err != nil {
			logger.Fatal(err)
		}
	},
}

// implement handlers.RecoveryHandlerLogger for zap logger
// type recoveryLogger struct {
// 	l *zap.SugaredLogger
// }

// func (p *recoveryLogger) Println(args ...any) {
// 	p.l.Error(args...)
// }
