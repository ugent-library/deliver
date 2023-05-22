package cli

import (
	"context"
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
	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstore"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/httpx/render"
	mw "github.com/ugent-library/middleware"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
	"github.com/urfave/cli/v2"
)

var appInfo = &struct {
	Branch string `json:"branch,omitempty"`
	Commit string `json:"commit,omitempty"`
}{
	Branch: os.Getenv("SOURCE_BRANCH"),
	Commit: os.Getenv("SOURCE_COMMIT"),
}

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

		// routes
		router.Get("/health", health.NewHandler(healthChecker))
		router.Get("/info", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, http.StatusOK, appInfo)
		})
		router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
			if err := hub.HandleWebSocket(w, r, r.URL.Query().Get("channels")); err != nil {
				logger.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		})
		router.Group(func(r *ich.Mux) {
			r.Use(
				func(next http.Handler) http.Handler {
					return http.MaxBytesHandler(next, config.MaxFileSize)
				},
				csrf.Protect(
					[]byte(config.Cookie.Secret),
					csrf.CookieName("deliver.csrf"),
					csrf.Path("/"),
					csrf.HttpOnly(true),
					csrf.SameSite(csrf.SameSiteStrictMode),
					csrf.FieldName("_csrf_token"),
				),
				// request context wrapper
				ctx.Set(ctx.Config{
					GetUserByRememberToken: repo.Users.GetByRememberToken,
					Router:                 router,
					ErrorHandlers: map[int]http.HandlerFunc{
						http.StatusNotFound:     errs.NotFound,
						http.StatusUnauthorized: errs.Unauthorized,
						http.StatusForbidden:    errs.Forbidden,
					},
					Permissions: permissions,
					Assets:      assets,
					Hub:         hub,
					Banner:      config.Banner,
				}),
			)

			// viewable by everyone
			r.NotFound(errs.NotFound)
			r.Get("/", pages.Home).Name("home")
			r.Get("/auth/callback", auth.Callback)
			r.Get("/login", auth.Login).Name("login")
			r.Get("/logout", auth.Logout).Name("logout")
			r.With(ctx.SetFolder(*repo.Folders)).Get("/share/{folderID}:{folderSlug}", folders.Share).Name("shareFolder")
			r.With(ctx.SetFile(*repo.Files)).Get("/files/{fileID}", files.Download).Name("downloadFile")
			// viewable by space owners and admins
			r.Group(func(r *ich.Mux) {
				r.Use(ctx.RequireUser)
				r.Get("/spaces", spaces.List).Name("spaces")
				r.With(ctx.RequireAdmin).Get("/new-space", spaces.New).Name("newSpace")
				r.With(ctx.RequireAdmin).Post("/spaces", spaces.Create).Name("createSpace")
				r.Route("/spaces/{spaceName}", func(r *ich.Mux) {
					r.Use(ctx.SetSpace(*repo.Spaces))
					r.Use(ctx.CanViewSpace)
					r.Get("/", spaces.Show).Name("space")
					r.Post("/folders", spaces.CreateFolder).Name("createFolder")
					r.With(ctx.RequireAdmin).Get("/edit", spaces.Edit).Name("editSpace")
					r.With(ctx.RequireAdmin).Put("/", spaces.Update).Name("updateSpace")
				})
				r.Route("/folders/{folderID}", func(r *ich.Mux) {
					r.Use(ctx.SetFolder(*repo.Folders))
					r.Use(ctx.CanEditFolder)
					r.Get("/", folders.Show).Name("folder")
					r.Get("/edit", folders.Edit).Name("editFolder")
					r.Put("/", folders.Update).Name("updateFolder")
					r.Post("/files", folders.UploadFile).Name("uploadFile")
					r.Delete("/", folders.Delete).Name("deleteFolder")
				})
				r.Route("/files/{fileID}", func(r *ich.Mux) {
					r.Use(ctx.SetFile(*repo.Files))
					r.Use(ctx.CanEditFile)
					r.Delete("/", files.Delete).Name("deleteFile")
				})
			})
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
