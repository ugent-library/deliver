package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oklog/ulid/v2"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	c "github.com/ugent-library/deliver/controllers"
	"github.com/ugent-library/deliver/crumb"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstore"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/middleware"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(appCmd)
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Start the web app server",
	Run: func(cmd *cobra.Command, args []string) {
		// setup services
		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			logger.Fatal(err)
		}

		storage, err := objectstore.New(config.Storage.Backend, config.Storage.Conn)
		if err != nil {
			logger.Fatal(err)
		}

		// setup permissions
		permissions := &models.Permissions{
			Admins: config.Admins,
		}

		// setup auth
		oidcAuth, err := oidc.NewAuth(context.TODO(), oidc.Config{
			URL:          config.OIDC.URL,
			ClientID:     config.OIDC.ID,
			ClientSecret: config.OIDC.Secret,
			RedirectURL:  config.OIDC.RedirectURL,
			CookieName:   "deliver.state",
			CookieSecret: []byte(config.Cookies.Secret),
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

		// htmx message hub
		hub := htmx.NewHub(htmx.Config{
			// TODO htmx secret config
			Secret: []byte(config.Cookies.Secret),
		})

		// controllers
		errs := c.NewErrorsController()
		auth := c.NewAuthController(repo, oidcAuth)
		pages := c.NewPagesController()
		spaces := c.NewSpacesController(repo)
		folders := c.NewFoldersController(repo, storage, config.MaxFileSize)
		files := c.NewFilesController(repo, storage)

		// request context wrapper
		wrap := c.Wrapper(c.Config{
			UserFunc:     repo.Users.GetByRememberToken,
			Router:       r,
			ErrorHandler: errs.HandleError,
			Permissions:  permissions,
			Assets:       assets,
			Hub:          hub,
			Banner:       config.Banner,
		})

		// router middleware
		r.Use(
			func(next http.Handler) http.Handler {
				return http.MaxBytesHandler(next, config.MaxFileSize)
			},
			crumb.Enable(
				crumb.WithErrorHandler(func(err error) {
					logger.Error(err)
				}),
			),
		)

		// routes
		r.NotFoundHandler = wrap(errs.NotFound)
		r.Handle("/", wrap(pages.Home)).Methods("GET").Name("home")
		r.Handle("/auth/callback", wrap(auth.Callback)).Methods("GET")
		r.Handle("/logout", wrap(auth.Logout)).Methods("GET").Name("logout")
		r.Handle("/login", wrap(auth.Login)).Methods("GET").Name("login")
		r.Handle("/spaces", wrap(c.RequireUser, spaces.List)).Methods("GET").Name("spaces")
		r.Handle("/spaces/{spaceName}", wrap(c.RequireUser, spaces.Show)).Methods("GET").Name("space")
		r.Handle("/new-space", wrap(c.RequireAdmin, spaces.New)).Methods("GET").Name("new_space")
		r.Handle("/spaces", wrap(c.RequireAdmin, spaces.Create)).Methods("POST").Name("create_space")
		r.Handle("/spaces/{spaceName}/edit", wrap(c.RequireAdmin, spaces.Edit)).Methods("GET").Name("edit_space")
		r.Handle("/spaces/{spaceName}", wrap(c.RequireAdmin, spaces.Update)).Methods("PUT").Name("update_space")
		r.Handle("/spaces/{spaceName}/folders", wrap(c.RequireUser, spaces.CreateFolder)).Methods("POST").Name("create_folder")
		r.Handle("/folders/{folderID}", wrap(c.RequireUser, folders.Show)).Methods("GET").Name("folder")
		r.Handle("/folders/{folderID}/edit", wrap(c.RequireUser, folders.Edit)).Methods("GET").Name("edit_folder")
		r.Handle("/folders/{folderID}", wrap(c.RequireUser, folders.Update)).Methods("PUT").Name("update_folder")
		r.Handle("/folders/{folderID}/files", wrap(c.RequireUser, folders.UploadFile)).Methods("POST").Name("upload_file")
		r.Handle("/folders/{folderID}", wrap(c.RequireUser, folders.Delete)).Methods("DELETE").Name("delete_folder")
		r.Handle("/files/{fileID}", wrap(files.Download)).Methods("GET").Name("download_file")
		r.Handle("/files/{fileID}", wrap(c.RequireUser, files.Delete)).Methods("DELETE").Name("delete_file")
		r.Handle("/share/{folderID}:{folderSlug}", wrap(folders.Share)).Methods("GET").Name("share_folder")

		mux := http.NewServeMux()
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			// TODO handle error
			hub.HandleWebSocket(w, r, r.URL.Query().Get("channels"))
		})
		mux.Handle("/", r)

		// apply these before request reaches the router
		handler := middleware.Apply(mux,
			middleware.Recover(func(err any) {
				if config.Production {
					logger.With(zap.Stack("stack")).Error(err)
				} else {
					logger.Error(err)
				}
			}),
			// apply before ProxyHeaders to avoid invalid referer errors
			csrf.Protect(
				[]byte(config.Cookies.Secret),
				csrf.CookieName("deliver.csrf"),
				csrf.Path("/"),
				csrf.Secure(config.Production),
				csrf.SameSite(csrf.SameSiteStrictMode),
				csrf.FieldName("_csrf_token"),
			),
			middleware.MethodOverride(
				middleware.MethodFromHeader(middleware.MethodHeader),
				middleware.MethodFromForm(middleware.MethodParam),
			),
			middleware.If(config.Production, handlers.ProxyHeaders),
			middleware.SetRequestID(func() string {
				return ulid.Make().String()
			}),
			zaphttp.SetLogger(logger.Desugar()),
			zaphttp.LogRequests(logger.Desugar()),
		)

		// start server
		// TODO make timeouts configurable
		addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		server := graceful.WithDefaults(&http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Minute,
			WriteTimeout: 10 * time.Minute,
		})
		logger.Infof("starting server at %s", addr)
		if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
			logger.Fatal(err)
		}
		logger.Info("gracefully stopped server")
	},
}
