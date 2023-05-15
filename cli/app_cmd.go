package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexliesenfeld/health"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/nics/ich"
	"github.com/ory/graceful"
	"github.com/ugent-library/deliver/controllers"
	"github.com/ugent-library/deliver/crumb"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstore"
	"github.com/ugent-library/deliver/repositories"
	mw "github.com/ugent-library/middleware"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
	"github.com/urfave/cli/v2"
)

var appCmd = &cli.Command{
	Name:  "app",
	Usage: "Start the web app server",
	Action: func(*cli.Context) error {
		// setup services
		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			return err
		}

		storage, err := objectstore.New(config.Storage.Backend, config.Storage.Conn)
		if err != nil {
			return err
		}

		// setup permissions
		permissions := &models.Permissions{
			Admins: config.Admins,
		}

		// setup auth
		oidcAuth, err := oidc.NewAuth(context.Background(), oidc.Config{
			URL:          config.OIDC.URL,
			ClientID:     config.OIDC.ID,
			ClientSecret: config.OIDC.Secret,
			RedirectURL:  config.OIDC.RedirectURL,
			CookieName:   "deliver.state",
			CookieSecret: []byte(config.Cookie.Secret),
		})
		if err != nil {
			return err
		}

		// setup health checker
		// TODO add checkers
		healthChecker := health.NewChecker()

		// setup assets
		assets, err := mix.New(mix.Config{
			ManifestFile: "static/mix-manifest.json",
			PublicPath:   "/static/",
		})
		if err != nil {
			return err
		}

		// setup router
		router := ich.New()
		router.Use(chimw.RequestID)
		if config.Production {
			router.Use(chimw.RealIP)
		}
		router.Use(mw.MethodOverride(
			mw.MethodFromHeader(mw.MethodHeader),
			mw.MethodFromForm(mw.MethodParam),
		))
		router.Use(zaphttp.SetLogger(logger.Desugar(), zapchi.RequestID))
		router.Use(chimw.RequestLogger(zapchi.LogFormatter()))
		router.Use(chimw.Recoverer)
		router.Use(chimw.StripSlashes)

		// htmx message hub
		hub := htmx.NewHub(htmx.Config{
			// TODO htmx secret config
			Secret: []byte(config.Cookie.Secret),
		})

		// controllers
		errs := controllers.NewErrorsController()
		auth := controllers.NewAuthController(repo, oidcAuth)
		pages := controllers.NewPagesController()
		spaces := controllers.NewSpacesController(repo)
		folders := controllers.NewFoldersController(repo, storage, config.MaxFileSize)
		files := controllers.NewFilesController(repo, storage)

		// request context wrapper
		wrap := controllers.Wrapper(controllers.Config{
			UserFunc:     repo.Users.GetByRememberToken,
			Router:       router,
			ErrorHandler: errs.HandleError,
			Permissions:  permissions,
			Assets:       assets,
			Hub:          hub,
			Banner:       config.Banner,
		})

		// routes
		router.Get("/health", health.NewHandler(healthChecker))
		// TODO clean this up, split off
		router.Get("/info", func(w http.ResponseWriter, r *http.Request) {
			info := &struct {
				Branch string `json:"branch,omitempty"`
				Commit string `json:"commit,omitempty"`
			}{
				Branch: os.Getenv("SOURCE_BRANCH"),
				Commit: os.Getenv("SOURCE_COMMIT"),
			}
			j, err := json.MarshalIndent(info, "", "  ")
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.Write(j)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		})
		router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
			// TODO handle error
			hub.HandleWebSocket(w, r, r.URL.Query().Get("channels"))
		})
		router.Get("/auth/callback", wrap(auth.Callback))
		router.Group(func(r *ich.Mux) {
			r.Use(
				func(next http.Handler) http.Handler {
					return http.MaxBytesHandler(next, config.MaxFileSize)
				},
				csrf.Protect(
					[]byte(config.Cookie.Secret),
					csrf.CookieName("deliver.csrf"),
					csrf.Path("/"),
					csrf.Secure(config.Production),
					csrf.SameSite(csrf.SameSiteStrictMode),
					csrf.FieldName("_csrf_token"),
				),
				crumb.Enable(
					crumb.WithErrorHandler(func(err error) {
						logger.Error(err)
					}),
				),
			)

			r.NotFound(wrap(errs.NotFound))
			r.Get("/", wrap(pages.Home)).Name("home")
			r.Get("/logout", wrap(auth.Logout)).Name("logout")
			r.Get("/login", wrap(auth.Login)).Name("login")
			r.Get("/spaces", wrap(controllers.RequireUser, spaces.List)).Name("spaces")
			r.Get("/spaces/{spaceName}", wrap(controllers.RequireUser, spaces.Show)).Name("space")
			r.Get("/new-space", wrap(controllers.RequireAdmin, spaces.New)).Name("newSpace")
			r.Post("/spaces", wrap(controllers.RequireAdmin, spaces.Create)).Name("createSpace")
			r.Get("/spaces/{spaceName}/edit", wrap(controllers.RequireAdmin, spaces.Edit)).Name("editSpace")
			r.Put("/spaces/{spaceName}", wrap(controllers.RequireAdmin, spaces.Update)).Name("updateSpace")
			r.Post("/spaces/{spaceName}/folders", wrap(controllers.RequireUser, spaces.CreateFolder)).Name("createFolder")
			r.Get("/folders/{folderID}", wrap(controllers.RequireUser, folders.Show)).Name("folder")
			r.Get("/folders/{folderID}/edit", wrap(controllers.RequireUser, folders.Edit)).Name("editFolder")
			r.Put("/folders/{folderID}", wrap(controllers.RequireUser, folders.Update)).Name("updateFolder")
			r.Post("/folders/{folderID}/files", wrap(controllers.RequireUser, folders.UploadFile)).Name("uploadFile")
			r.Delete("/folders/{folderID}", wrap(controllers.RequireUser, folders.Delete)).Name("deleteFolder")
			r.Get("/files/{fileID}", wrap(files.Download)).Name("downloadFile")
			r.Delete("/files/{fileID}", wrap(controllers.RequireUser, files.Delete)).Name("deleteFile")
			r.Get("/share/{folderID}:{folderSlug}", wrap(folders.Share)).Name("shareFolder")
		})

		// start server
		// TODO make timeouts configurable
		addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		server := graceful.WithDefaults(&http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  10 * time.Minute,
			WriteTimeout: 10 * time.Minute,
		})
		logger.Infof("starting server at %s", addr)
		if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
			return err
		}
		logger.Info("gracefully stopped server")

		return nil
	},
}
